# Lottery-driven verification approach - evaluation

## Costs

**TODO**

most likely the costly parts are:

1) cost of deposits placed by Providers
2) gas costs (unless task execution done off-Ethereum)
    1) Subtask assignment and reshuffling
    2) atomic swap scheme, incl. punishing non-revealing Providers
    3) rng - if more involved than suggested
    4) commiting to hashes of results, punishing incorrect commitments
    5) commiting to task division, punishing incorrect commitments

## Guarantees and risks for parties

### Notation

![alt text](https://latex.codecogs.com/gif.latex?K) - number of Providers applying for a Task
![alt text](https://latex.codecogs.com/gif.latex?K_T) - number of Providers elected to handle Task
![alt text](https://latex.codecogs.com/gif.latex?K^*) - number of Sybil identities (addresses) acting as Providers (applying for a Task)
![alt text](https://latex.codecogs.com/gif.latex?P_{pow}) - computing power of (some) Provider
![alt text](https://latex.codecogs.com/gif.latex?T_{pow}) - computing power of all Provider's elected for Task
![alt text](https://latex.codecogs.com/gif.latex?R_{pow}) - computing power of Requestor
![alt text](https://latex.codecogs.com/gif.latex?F) - fee for Task
![alt text](https://latex.codecogs.com/gif.latex?C) - cost of computing a Task

Guarantees for Requestors:

1. Assuming no Sybil attack, Provider never gets paid for bad work: Provider that does 100% incorrect work (delivers junk), has 100% chance of getting 0 reward
2. Redundant work is done meticulously, because it is done by Requestor itself
    - to add to that, there is room for having incentive for Providers to do redundant work meticulously and catch cheaters

Guarantees for Providers:

1. Assuming no Sybil attack, Provider that does 100% correct work, still gets his expected payout in the long run

Risks for Requestors:

#### Opportunistic Sybils of P

**Adding Sybil addresses to a Task allows Provider to introduce incorrect results and increase expected payout**

P1 puts his Sybil identities P\* on a Task. P1 provides good results, P\* push empty results cheaply. If they are caught, P1 gets some more work to do. If P1 gets the winning-ticket, P1 runs away with payout leaving R with bad results.

    - unfortunately this leaves P1 better off (winning strategy for all Providers)
    - With P\* being there the expected reward per subtask is raised from ![alt text](https://latex.codecogs.com/gif.latex?\frac{1}{N}) to ![alt text](https://latex.codecogs.com/gif.latex?\frac{N+S^*}{N^2}), where ![alt text](https://latex.codecogs.com/gif.latex?N) - number of subtask, ![alt text](https://latex.codecogs.com/gif.latex?S^*) number of empty results and there's a single winning ticket. (for by derivation of this by ŁG see [here](https://github.com/imapp-pl/golem_rd/blob/4c060e48978dfd593809d54ab44a34d89649e036/Lottery-driven-verification.md), known problems section)
    
From a more general perspective, consider 2 strategies for P1: 
  - **"honest"** where P1 holds only one address (![alt text](https://latex.codecogs.com/gif.latex?K^*=1) and gets to calculate a Task
  - **"cheating"** where P1 holds two addresses (![alt text](https://latex.codecogs.com/gif.latex?K^*=2) and gets to calculate a Task with one or two of these (assuming Providers are picked to the Task randomly: (![alt text](https://latex.codecogs.com/gif.latex?K^T) out of (![alt text](https://latex.codecogs.com/gif.latex?K) total applying.
  
The gain from employing the "cheating" strategy compared to the "honest" one, for a single Task is (derivation skipped):

(![alt text](https://latex.codecogs.com/gif.latex?E(\text{payout}|\text{cheating})-E(\text{payout}|\text{honest})=2\cdot\frac{P_{pow}}{T_{pow}}\frac{(K_T-1)(K-K_T)}{(K-1)(K-2)}\cdot C \cdot \frac{K_T-1}{K_T^2})

Intuitively, this value is maxed out for a "right" proportion of (![alt text](https://latex.codecogs.com/gif.latex?K_T) within (![alt text](https://latex.codecogs.com/gif.latex?K), and generally decreases with (![alt text](https://latex.codecogs.com/gif.latex?K) increasing. 
The more Providers apply for a single Task, the better.

An example value of such gain is $0.0015 for a Task costing $1 to compute, assuming P1 has 10% of compute power and his 2 Sybils are one of 100 Providers applying (![alt text](https://latex.codecogs.com/gif.latex?K) for a 10-Provider (![alt text](https://latex.codecogs.com/gif.latex?K_T=10)) Task.

Other
1. Task might take long to compute or fail entirely due to uncertanities in difficulty/pricing
3. Random number generator might be too lame to provide honest lotteries and subtask assignment
4. Requestor risks having undetected incorrect results introduced by **malicious** Provider (griefing)

Risks for Providers:

1. "no reward till bored" is probable for small providers
1. Adding Sybil (Provider) addresses to a Task allows **Requestor** to be charged for computation less than what honest Providers expect
3. Provider risks sacrificing some work for sake of determining Task difficulty and pricing
4. Provider risks being dropped and his work rejected by **malicious** Requestor reporting an incorrect result unjustly (griefing)

## External constraints and implications

1. Determinism of results required
2. Requestor needs to provide "best guess" of Task difficulty/price in advance and know his maximum price he's willing to pay
3. Tasks need to be finite and divisible into chunks (subtasks)
4. High-latency protocol

## Implementation risks

1. Many pieces of the protocol still to be fleshed out:
    - commitment to division of a broadcast task by Requestor
    - what happens if KDF reveal phase breaks
    - reshuffling details are tricky
    - consider market properties of the price-bumping mechainsm - is it exploitable? is it practical enough?
    - careful random number generator considerations
2. Subtask assignment and reshuffling might prove to be to unwieldy to be handled properly and cheaply enough. This might require introducing of some additional BFT mechanism (like e.g. deployment of an ad-hoc Tendermint sidechain for single Task execution).

## Extensibility

**TODO**

1. this compatible with Transaction Framework Model idea?
