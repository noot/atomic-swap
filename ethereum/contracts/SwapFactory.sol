// SPDX-License-Identifier: LGPLv3
pragma solidity ^0.8.5;

import "./Secp256k1.sol";

contract SwapFactory is Secp256k1 {

    // Swap state is PENDING when the swap is first created and funded
    // Alice sets Stage to READY when she sees the funds locked on the other chain.
    // this prevents Bob from withdrawing funds without locking funds on the other chain first
    // Stage is set to COMPLETED upon the swap value being claimed or refunded.

    enum Stage {
        INVALID,
        PENDING,
        READY,
        COMPLETED
    }

    struct Swap {
        // contract creator, Alice
        address payable owner;

        // address allowed to claim the ether in this contract
        address payable claimer;

        // the keccak256 hash of the expected public key derived from the secret `s_b`.
        // this public key is a point on the secp256k1 curve
        bytes32 pubKeyClaim;

        // the keccak256 hash of the expected public key derived from the secret `s_a`.
        // this public key is a point on the secp256k1 curve
        bytes32 pubKeyRefund;

        // timestamp (set at contract creation)
        // before which Alice can call either set_ready or refund
        uint256 timeout_0;

        // timestamp after which Bob cannot claim, only Alice can refund.
        uint256 timeout_1;

        // the value of this swap.
        uint256 value;

        // choose random
        uint256 nonce;
    }

    mapping(bytes32 => Stage) public swaps;

    event New(bytes32 swapID, bytes32 claimKey, bytes32 refundKey, uint256 timeout_0, uint256 timeout_1);
    event Ready(bytes32 swapID);
    event Claimed(bytes32 swapID, bytes32 s);
    event Refunded(bytes32 swapID, bytes32 s);

    // new_swap creates a new Swap instance with the given parameters.
    // it returns the swap's ID.
    function new_swap(bytes32 _pubKeyClaim, 
        bytes32 _pubKeyRefund, 
        address payable _claimer, 
        uint256 _timeoutDuration,
        uint256 _nonce
    ) public payable returns (bytes32) {

        Swap memory swap;
        swap.owner = payable(msg.sender); 
        swap.claimer = _claimer;
        swap.pubKeyClaim = _pubKeyClaim;
        swap.pubKeyRefund = _pubKeyRefund;
        swap.timeout_0 = block.timestamp + _timeoutDuration;
        swap.timeout_1 = block.timestamp + (_timeoutDuration * 2);
        swap.value = msg.value;
        swap.nonce = _nonce;

        bytes32 swapID = keccak256(abi.encode(swap));

        // make sure this isn't overriding an existing swap
        require(swaps[swapID] == Stage.INVALID);

        emit New(swapID, _pubKeyClaim, _pubKeyRefund, swap.timeout_0, swap.timeout_1);
        swaps[swapID] = Stage.PENDING;
        return swapID;
    }

    // Alice must call set_ready() within t_0 once she verifies the XMR has been locked
    function set_ready(Swap memory _swap) public {
        bytes32 swapID = keccak256(abi.encode(_swap));
        require(swaps[swapID] == Stage.PENDING, "swap is not in PENDING state");
        require(_swap.owner == msg.sender, "only the swap owner can call set_ready");
        swaps[swapID] = Stage.READY;
        emit Ready(swapID);
    }

    // is_ready returns whether a swap has been set to "ready" or not.
    // note: it will return false, not revert, if the swap does not exist.
    function is_ready(bytes32 _swapID) public view returns (bool) {
        return swaps[_swapID] == Stage.READY;
    }

    // Bob can claim if:
    // - Alice doesn't call set_ready or refund within t_0, or
    // - Alice calls ready within t_0, in which case Bob can call claim until t_1
    function claim(Swap memory _swap, bytes32 _s) public {
        bytes32 swapID = keccak256(abi.encode(_swap));
        Stage swapStage = swaps[swapID];
        require(swapStage != Stage.COMPLETED && swapStage != Stage.INVALID, "swap is already completed");
        require(msg.sender == _swap.claimer, "only claimer can claim!");
        require((block.timestamp >= _swap.timeout_0 || swapStage == Stage.READY), "too early to claim!");
        require(block.timestamp < _swap.timeout_1, "too late to claim!");

        verifySecret(_s, _swap.pubKeyClaim);
        emit Claimed(swapID, _s);

        // send eth to caller (Bob)
        _swap.claimer.transfer(_swap.value);
        swaps[swapID] = Stage.COMPLETED;
    }

    // Alice can claim a refund:
    // - Until t_0 unless she calls set_ready
    // - After t_1, if she called set_ready
    function refund(Swap memory _swap, bytes32 _s) public {
        bytes32 swapID = keccak256(abi.encode(_swap));
        Stage swapStage = swaps[swapID];
        require(swapStage != Stage.COMPLETED && swapStage != Stage.INVALID, "swap is already completed");
        require(msg.sender == _swap.owner, "refund must be called by the swap owner");
        require(
            block.timestamp >= _swap.timeout_1 ||
            (block.timestamp < _swap.timeout_0 && swapStage != Stage.READY),
            "it's the counterparty's turn, unable to refund, try again later"
        );

        verifySecret(_s, _swap.pubKeyRefund);
        emit Refunded(swapID, _s);

        // send eth back to owner==caller (Alice)
        _swap.owner.transfer(_swap.value);
        swaps[swapID] = Stage.COMPLETED;
    }

    function verifySecret(bytes32 _s, bytes32 pubKey) internal pure {
        require(
            mulVerify(uint256(_s), uint256(pubKey)),
            "provided secret does not match the expected public key"
        );
    }
}
