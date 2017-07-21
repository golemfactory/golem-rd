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

### Guarantees for Requestors:

1. Assuming no Sybil attack, Provider never gets paid for bad work: Provider that does 100% incorrect work (delivers junk), has 100% chance of getting 0 reward
2. Redundant work is done meticulously, because it is done by Requestor itself
    - to add to that, there is room for having incentive for Providers to do redundant work meticulously and catch cheaters

### Guarantees for Providers:

1. Assuming no Sybil attack, Provider that does 100% correct work, still gets his expected payout in the long run

### Risks for Requestors:

#### Opportunistic Sybils of P

**Adding Sybil addresses to a Task allows Provider to introduce incorrect results and increase expected payout**

P1 puts his Sybil identities P\* on a Task. P1 provides good results, P\* push empty results cheaply. If they are caught, P1 gets some more work to do. If P1 gets the winning-ticket, P1 runs away with payout leaving R with bad results.

- unfortunately this leaves P1 better off (winning strategy for all Providers)
- With P\* being there the expected reward per subtask is raised from ![alt text](https://latex.codecogs.com/gif.latex?\frac{1}{N}) to ![alt text](https://latex.codecogs.com/gif.latex?\frac{N+S^*}{N^2}), where ![alt text](https://latex.codecogs.com/gif.latex?N) - number of subtask, ![alt text](https://latex.codecogs.com/gif.latex?S^*) number of empty results and there's a single winning ticket. (for by derivation of this by ŁG see [here](https://github.com/imapp-pl/golem_rd/blob/4c060e48978dfd593809d54ab44a34d89649e036/Lottery-driven-verification.md), known problems section)
    
From a more general perspective, consider 2 strategies for P1: 
  - **"honest"** where P1 holds only one address ![alt text](https://latex.codecogs.com/gif.latex?K^*=1) and gets to calculate a Task
  - **"cheating"** where P1 holds two addresses ![alt text](https://latex.codecogs.com/gif.latex?K^*=2) and gets to calculate a Task with one or two of these (assuming Providers are picked to the Task randomly: ![alt text](https://latex.codecogs.com/gif.latex?K^T) out of ![alt text](https://latex.codecogs.com/gif.latex?K) total applying.)
  
The gain from employing the "cheating" strategy compared to the "honest" one, for a single Task is (derivation skipped):

![alt text](https://latex.codecogs.com/gif.latex?E(\text{payout}|\text{cheating})-E(\text{payout}|\text{honest})=)

![alt text](https://latex.codecogs.com/gif.latex?=2\cdot\frac{P_{pow}}{T_{pow}}\frac{(K_T-1)(K-K_T)}{(K-1)(K-2)}\cdot&space;C\cdot\frac{K_T-1}{K_T^2})

Intuitively, this value is maxed out for some proportion of ![alt text](https://latex.codecogs.com/gif.latex?K_T) within ![alt text](https://latex.codecogs.com/gif.latex?K), and generally decreases with ![alt text](https://latex.codecogs.com/gif.latex?K) increasing. 
The more Providers apply for a single Task, the better.

An example value of such gain is $0.0015 for a Task costing $1 to compute, assuming P1 has 10% of compute power and his 2 Sybils are one of 100 Providers applying (![alt text](https://latex.codecogs.com/gif.latex?K=100)) for a 10-Provider (![alt text](https://latex.codecogs.com/gif.latex?K_T=10)) Task.

#### Difficulty vs pricing

**Task might take long to compute or fail entirely due to uncertanities in difficulty/pricing**

details TODO

#### RNG

**Random number generator might be too lame to provide honest lotteries and subtask assignment**

#### Griefing by incorrect results

**Requestor risks having undetected incorrect results introduced by **malicious** Provider (griefing)**

### Risks for Providers:

#### No reward till bored

The risks consist in the possibility of Provider never winning any lottery before some critical "trial time" passes.
Although, the expected value of income checks out and is "fair", we still worry about the possibility of actual earnings of a Provider to be dissapointing.
Let's model this dissapointment using the two critical values:

![alt text](https://latex.codecogs.com/gif.latex?top_1) - ("time-of-patience") - if no income for Provider within that time, Provider is dissapointed with Golem.

![alt text](https://latex.codecogs.com/gif.latex?top_2) - if less than 90% expected income for Provider within that time, Provider is dissapointed with Golem

Additionally define: ![alt text](https://latex.codecogs.com/gif.latex?T) - time it takes to calculate the Task, ![alt text](https://latex.codecogs.com/gif.latex?F(n,p,k)) - CDF of ![alt text](https://latex.codecogs.com/gif.latex?Binomial(n,p)), i.e. probability that at most ![alt text](https://latex.codecogs.com/gif.latex?k) out of ![alt text](https://latex.codecogs.com/gif.latex?n) succeed, with probability of single success ![alt text](https://latex.codecogs.com/gif.latex?p).

Then probability to be dissapointed for the two above reasons is respectively:

![alt text](https://latex.codecogs.com/gif.latex?P_1=F\left(\frac{top_1\cdot&space;T_{pow}}{T\cdot&space;R_{pow}},\frac{P_{pow}}{T_{pow}},0\right))

![alt text](https://latex.codecogs.com/gif.latex?P_2=F\left(\frac{top_2\cdot&space;T_{pow}}{T\cdot&space;R_{pow}},\frac{P_{pow}}{T_{pow}},0.9\cdot&space;\frac{top_2\cdot&space;P_{pow}}{T\cdot&space;R_{pow}}\right))

Note that the number of lotteries taken within the "time-of-patience" is ![alt text](https://latex.codecogs.com/gif.latex?\frac{top\cdot&space;T_{pow}}{T\cdot&space;R_{pow}}), while probability of winnig a lottery for a Task is proportional to the proportion of computing power of Provider within the pool of computing power ammassed for a single Task: ![alt text](https://latex.codecogs.com/gif.latex?\frac{P_{pow}}{T_{pow}}).

These formulas yield an example estimation. Assume:
 - a Task requires ![alt text](https://latex.codecogs.com/gif.latex?T=\text{10&space;hours}) of Requestors CPU to compute
 - Provider expects to earn something overnight ![alt text](https://latex.codecogs.com/gif.latex?top_1=\text{12&space;hours}) 
 - get 90% of expected income over a week ![alt text](https://latex.codecogs.com/gif.latex?top_2=\text{84&space;hours}) (seven overnight sessions)
 - Provider's and Requestor's power is the same and about 40% of whole Tasks computing power.
 
These numbers arise when one considers getting a 10 hour task done in 4 hours and Requestor having a comparable machine to Provider's.
(There is an assumption made that we worry most about early adopters).

Under such assumptions ![alt text](https://latex.codecogs.com/gif.latex?P_1\approx0.216) and ![alt text](https://latex.codecogs.com/gif.latex?P_2\approx0.350), which are quite worrying.

The above results take many general assumptions and cut some corners like: all Tasks are equal, all Tasks complete, very fine division into subtasks, single winning ticket per Task etc.

#### Opportunistic P-Sybils of R

**Adding Sybil (Provider) addresses to a Task allows **Requestor** to be charged for computation less than what honest Providers expect**

R can put some own P\*s on his task and have them _never calculate anything_, but commit to some fake hashes and grab the winning tickets, if P\*s get a lot of them. It may turn out, that all winning tickets go to P\*s (hence back to R), and none to working Ps. Then R resubmits the part of the Task which was not calculated as a new Task2. This strategy has better expected return for R, if only Task submitting and Sybil identities are relatively cheap.

Consider 2 strategies for R: 
  - **"honest"** where R holds no Provider-addresses ![alt text](https://latex.codecogs.com/gif.latex?K^*=0)
  - **"cheating"** where R holds Provider-addresses ![alt text](https://latex.codecogs.com/gif.latex?K^*>=1) and may get to calculate R's own Task with some of them (assuming Providers are picked to the Task randomly: ![alt text](https://latex.codecogs.com/gif.latex?K^T) out of ![alt text](https://latex.codecogs.com/gif.latex?K) total applying.)
  
The gain from employing the "cheating" strategy compared to the "honest" one, for a single Task is (derivation skipped):

![alt text](https://latex.codecogs.com/gif.latex?E(\text{cost}|\text{honest})-E(\text{cost}|\text{cheating})=)

![alt text](https://latex.codecogs.com/gif.latex?=F\cdot\frac{K_T-1}{K_T^2}\cdot\sum_{d=0}^{K^*}dP(\text{exactly&space;d&space;Sybil&space;Providers&space;got&space;in}))

(**TODO**: simplify the above formula, it seems it is linear in ![alt text](https://latex.codecogs.com/gif.latex?K^*))

Suppose Requestor has ![alt text](https://latex.codecogs.com/gif.latex?K^*) Sybil Provider address out of 100 Providers applying (![alt text](https://latex.codecogs.com/gif.latex?K=100)) for a 10-Provider (![alt text](https://latex.codecogs.com/gif.latex?K_T=10)) Task.
Then his expected gain is $0.009 for a Task with fee equal $1, for every Sybil Provider address R has.

This in fact is huge savings opportunity for R and ought to be addressed!

#### Wasted effort for pricing

**Provider risks sacrificing some work for sake of determining Task difficulty and pricing**

#### Griefing by unjust rejection

**Provider risks being dropped and his work rejected by **malicious** Requestor reporting an incorrect result unjustly (griefing)**

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
