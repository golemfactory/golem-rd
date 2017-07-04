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

Guarantees for Requestors:

1. Assuming no Sybil attack, Provider never gets paid for bad work: Provider that does 100% incorrect work (delivers junk), has 100% chance of getting 0 reward
2. Redundant work is done meticulously, because it is done by Requestor itself
    - to add to that, there is room for having incentive for Providers to do redundant work meticulously and catch cheaters

Guarantees for Providers:

1. Assuming no Sybil attack, Provider that does 100% correct work, still gets his expected payout in the long run

Risks for Requestors:

1. Adding Sybil addresses to a Task allows **Provider** to introduce incorrect results and increase expected payout
2. Task might take long to compute or fail entirely due to uncertanities in difficulty/pricing
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

## Implementation risks

1. Many pieces of the protocol still to be fleshed out:
    - commitment to division of a broadcast task by Requestor
    - what happens if KDF reveal phase breaks
    - reshuffling details are tricky
    - consider market properties of the price-bumping mechainsm - is it exploitable? is it practical enough?
    - careful random number generator considerations
2. Subtask assignment and reshuffling might prove to be to unwieldy to be handled properly and cheaply enough. This might require introducing of some additional BFT mechanism (like e.g. deployment of an ad-hoc Tendermint sidechain for single Task execution).
3. Protocol is exploitable if Sybil identities are cheap

## Extensibility

**TODO**

1. this compatible with Transaction Framework Model idea?