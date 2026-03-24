// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract AletheiaProof {
    struct Proof {
        string sha256;
        string ipfsCID;
        string merkleRoot;
        uint256 timestamp;
        address owner;
    }

    mapping(string => Proof) private proofs;
    mapping(string => bool) private exists;

    event ProofStored(string indexed sha256, string ipfsCID, string merkleRoot, uint256 timestamp, address indexed owner);

    function storeProof(string memory _sha256, string memory _ipfsCID, string memory _merkleRoot) public {
        require(!exists[_sha256], "Proof already exists for this hash");

        proofs[_sha256] = Proof({
            sha256: _sha256,
            ipfsCID: _ipfsCID,
            merkleRoot: _merkleRoot,
            timestamp: block.timestamp,
            owner: msg.sender
        });

        exists[_sha256] = true;

        emit ProofStored(_sha256, _ipfsCID, _merkleRoot, block.timestamp, msg.sender);
    }

    function getProof(string memory _sha256) public view returns (string memory, string memory, string memory, uint256, address) {
        require(exists[_sha256], "Proof does not exist");
        Proof memory p = proofs[_sha256];
        return (p.sha256, p.ipfsCID, p.merkleRoot, p.timestamp, p.owner);
    }
}
