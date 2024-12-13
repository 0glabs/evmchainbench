# Performance Analysis of 0G, Evmos, Kava, Bera, and Sei

| Chain   | Simple | ERC20 | Uniswap |
|---------|--------|-------|---------|
| 0G      | 769    | 369   | 164     |
| Bera    | 910    | 638   | 224     |
| Evmos   | 790    | 859   | 689     |
| Kava    | 637    | 84    | 36      |
| Sei     | 784    | 784   | 392     |

## Key Observations and Insights

### 1. Cosmos+Ethermint vs. Cosmos+Beacon API+Geth/Reth
- **0G, Evmos, and Kava** use Cosmos+Ethermint, where each Ethereum transaction is wrapped into a Cosmos transaction for consensus processing. This introduces additional overhead compared to directly processing Ethereum transactions in an EVM.
- **Bera** adopts Cosmos+Beacon API+Geth/Reth, wrapping an entire Ethereum block payload into a single Cosmos transaction. This significantly reduces the transaction load on the consensus layer, resulting in better performance across all test categories.

### 2. Sei's Unique Modifications
- **Sei** has extensively modified Cosmos, Tendermint, and Go-Ethereum. These deep changes make it fundamentally different from standard Cosmos chains like 0G, Evmos, and Kava. As such, Sei's performance cannot be directly compared to other chains in this analysis.

### 3. Block Production in Cosmos+Ethermint
- Ethermint-based chains produce blocks based on Ethereum transactions' **gas limits** rather than **gas used**, as the Cosmos consensus layer cannot calculate gas used during block production. To prevent misuse of inflated gas limits, a **minimum gas usage ratio** (typically 50%, adjustable via consensus) is enforced.
- **Evmos’ higher TPS** in ERC20 and Uniswap tests is primarily due to its larger block size configuration, not inherent performance optimizations.

### 4. Performance Gap Between 0G and Kava
- **0G and Kava** share similar block size configurations, but **0G achieves better TPS** due to an improved `estimateGas` method. This enhancement allows for more accurate gas limit calculations, improving transaction processing efficiency.

### 5. Bera's Design Advantage
- Bera's use of block-level payload processing avoids the overhead of transaction-by-transaction consensus inherent in Ethermint chains. This architectural decision gives Bera a significant performance edge in all test scenarios.

## Conclusion
The performance differences highlight the impact of architectural and implementation choices:
- **Bera** excels due to its block-level payload processing approach, which reduces consensus overhead.
- **0G** demonstrates the advantages of refining critical components like `estimateGas` to improve transaction throughput.
- **Evmos** achieves high ERC20 and Uniswap TPS through increased block size, while **Kava** lags due to older CometBFT and less efficient gas estimation.
- **Sei’s extensive customizations** set it apart from other chains, making direct comparisons to standard Cosmos-based architectures inappropriate.
