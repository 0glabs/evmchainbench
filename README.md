# Performance Analysis of 0G, Evmos, Kava, Bera, and Sei

| **Chain** | **Simple TPS** | **ERC20 TPS** | **Uniswap TPS** |
|-----------|----------------|---------------|-----------------|
| **0G**    | 769            | 369           | 164             |
| **Bera**  | 910            | 638           | 224             |
| **Evmos** | 790            | 859           | 689             |
| **Kava**  | 637            | 84            | 36              |
| **Sei**   | 784            | 784           | 392             |

### Key Clarifications

1. **TPS Measurements**  
   The numbers shown in the chart represent **transactions per second (TPS)** achieved during our benchmarking. These results are based on specific configurations we tested and optimized to ensure fairness. While we aimed to make the comparison competitive by adjusting configurations for each chain, we acknowledge that some configurations may not be fully optimal, as optimizing every chain requires significant time investment. 

2. **Benchmarking Methodology**  
   This repository enables **anyone** to benchmark these chains using the same testing suite we used, ensuring **truthful and reproducible results**. To illustrate this, we recorded a demo showcasing how we achieved these performance numbers. This demo is particularly helpful for non-technical users or those who prefer a quick overview instead of running the tests themselves.

---

## Key Observations and Insights

### 1. Cosmos+Ethermint vs. Cosmos+Beacon API+Geth/Reth  
- **0G, Evmos, and Kava** use Cosmos+Ethermint, where each Ethereum transaction is wrapped into a Cosmos transaction for consensus processing. This approach introduces additional consensus overhead, resulting in slightly lower TPS compared to Bera.
- **Bera** leverages a **Cosmos+Beacon API+Geth/Reth** design, which wraps entire Ethereum block payloads into single Cosmos transactions. This block-level payload processing significantly reduces consensus overhead, yielding consistently higher TPS across all categories.

### 2. Why Include Sei?  
While Sei's extensive modifications to Cosmos, Tendermint, and Go-Ethereum make its architecture fundamentally different, we included Sei here for completeness. However, given Sei’s deviation from the standard Cosmos framework, direct comparisons may not fully reflect differences in architectural efficiency.

### 3. Evmos’ Higher TPS and Block Configuration Tradeoffs  
Evmos achieves higher ERC20 and Uniswap TPS primarily due to its **larger block size** configuration, rather than inherent performance optimizations.  
- **Tradeoffs of larger block size**:  
  - **Finality**: Larger blocks increase finality time, delaying confirmation of transactions.  
  - **Latency**: They can introduce higher network latency due to increased transmission times.  
  - **Network Congestion**: Larger blocks may cause higher peak loads, leading to temporary congestion under high transaction volume.  
  - **Gas Fees**: Larger blocks can reduce gas fees during normal conditions but may increase volatility during congestion.  

In **v2**, we are carefully balancing these tradeoffs by mimicking Bera’s execution client approach while **reducing block finality time** to achieve both high TPS and low-latency finality.

### 4. Performance Gap Between 0G and Kava  
Despite similar block size configurations, **0G outperforms Kava** due to its enhanced `estimateGas` method. This improvement allows for more precise gas limit calculations, reducing unnecessary overhead and boosting overall transaction throughput.

### 5. Bera’s Block-Level Payload Advantage  
Bera’s architecture avoids the transaction-by-transaction consensus overhead of Ethermint chains by processing entire blocks as single payloads. This design gives Bera a clear performance edge in all test categories.

---

## Next Steps: V2 Enhancements  
In our upcoming **v2 release**, we are:  
- **Adopting a block-level payload approach** similar to Bera’s design to reduce consensus overhead.  
- **Optimizing block finality time** to ensure fast transaction confirmation without sacrificing high TPS.  
- **Conducting further benchmarks** to refine performance tradeoffs between block size, latency, and gas fees.
