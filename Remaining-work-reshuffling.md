# Remaining work reshuffling

We propose an extension of the scheme outlined in [[Work amount measurement by fair sampling ex ante]], where providers beat each other in a race to be assigned new subtasks.
The problem identified in this scheme is that it enables attacks, as requestor has final say on which provider is dealt which subtask - even if the task ordering is randomly predetermined and committed to. 
Previously suggested mitigation was to apply BFT replicated state which prohibits manipulating of the assignment of the subtasks, e.g. requestor and providers forming a disposable, ad-hoc Tendermint network. 
This is potentially expensive and complex.

The proposition "Remaining work reshuffling" is an extension of this race protocol, which 
 - sustains the motivation to race (and thus the accuracy of ex-ante task pricing), 
 - prohibits the manipulation of subtask assignment, and 
 - doesn't require BFT replicated state aside the Ethereum chain smart-contracts.

It is also believed that it doesn't increase smart-contract activity _too much_ (to be analyzed).

## Protocol outline

Computation of a single task proceeds as follows:

1. Requestor defines the task as usual, defines the division of task into subtasks. 
Requestor commits to that on-chain (e.g. publishes a merkle root of the inputs). 
Violating this commitment results requestor being slashed.
2. The commitment on chain fixes seed `s1` which is determined by some non-foreseeable and non-manipulable data in the near future.
The assumption is that `s1` is in no way under control of requestor or providers.
3. Providers `P1`, ..., `Pn` sign up for the task
4. `s1` shuffles the subtasks and determines assignment matrix `M1`, which assigns equal count of subtasks to `P1`, ..., `Pn` (**NOTE** order matters!), e.g.:

    |           | P1 | P2 | P3 | P4 |
    |-----------|----|----|----|----|
    | subtasks: | 4  | 11 | 12 | 1  |
    |           | 5  | 6  | 7  | 15 |
    |           | 2  | 8  | 14 | 16 |
    |           | 10 | 9  | 3  | 13 |

    (Assuming 4 providers and 16 subtasks)
5. `P1`, ..., `Pn` compute the tasks as fast as possible (motivation argued below).
6. First to complete his pool of subtasks, publishes their results' merkle root on-chain (can be done off-chain too, see below) and by that starts a challenge.
Challenge means other providers must likewise publish their roots, committing to their obtained results.
Those other providers have limited time for this.
Note here that only the results of given provider's _first-most_ subtasks can be published and are payable -- order matters.
7. All subtasks that have not been published this way are brought into a common pool, reshuffled using `s2` (same properties as `s1`), and assigned using matrix `M2`.
If after first challenged `P1`, ..., `P4` published the results for 1, 4, 1, 2 of their subtasks respectively, i.e.:

    |           | P1 | P2 | P3 | P4 |
    |-----------|----|----|----|----|
    | results:  | 4  | 11 | 12 | 1  |
    |           |    | 6  |    | 15 |
    |           |    | 8  |    |    |
    |           |    | 9  |    |    |

    Then the `M2` is:

    |           | P1 | P2 | P3 | P4 |
    |-----------|----|----|----|----|
    | subtasks: | 7  | 5  | 16 | 2  |
    |           | 3  | 10 | 3  | 13 |

    Eight subtasks carried on to round 2, reshuffled and distributed uniformly.

8. Lather, rinse, repeat

## Calculation of the price

### Price per hour

Calculation of the payout proceeds the same way as in the aforementioned solution, i.e. average subtask time is sampled out of all times ex-ante. This, together with number of completed subtasks per provider and an agreed price per hour determine the final payout.

### Price per subtask

Modification proposed by Grzegorz: no price per hour but price per subtask.
The mechanism of reshuffling the remainder of subtask is still useful, as it provides non-repudiable and non-controllable ordering and assignment of tasks.

Pros of "price per subtask":
  - lesser vector of attacks (no "Sybil-ing providers" attack, see below)
  - simpler protocol, less analysis required
Cons:
  - market less effective. Price varies from task to task so requestor does trial and error when setting task price.
  - providers might abandon difficult under-priced subtasks
  - tool to estimate subtask price required

## Shuffling mechanism rationale

The properties of the reshuffling -- i.e. the assignment and order of subtasks such that is not-foreseeable and non-manipulable -- prohibits any party from ordering the subtasks in a way that could drive the price up or down.
The assignment and order of subtasks once committed, cannot be repudiated.
In order to be paid, the provider must calculate subtasks as prescribed.

In particular, the requestor cannot flood the task with fake tasks (instantaneous to compute, drive the average cost down) and deal those to fake providers in his control.
If that was attempted, the fake tasks would spread evenly between all of the providers, including honest ones.

Similarly, the requestor cannot exploit his information about the difficulty of subtasks and deal the cheapest ones to providers in his control.

## Challenge every round drives the race

In an ideal setting, where both `P1`, ..., `Pn` have same capacity *and* the subtasks are uniformly difficult, `P1`, ..., `Pn` all finish work in same moment. 
Whichever publishes the challenge, triggers the remainder to publish too, and the task is completed -- no need for reshuffling.

Provider is motivated to do the challenge because it increases his chance of getting more subtasks in next rounds and maximizes his payout.
That's because he expects to get unstarted tasks of the laggards in subsequent rounds.
Since this motivation, he will not try to delay the publishing to pretend that tasks where more difficult than they really were (to pump the final payout for everyone).
These motivations are considered further in next sections.

Note that in the real-world, there will be winners and laggards in the race since tasks are unevenly difficult and computing capacity varies.
Regardless, the mechanism provides strong guarantees that they will be paid for all the work they've done.
There are two variants of the protocol worth mentioning: on reshuffling, the non-finished-but-on-the-top-of-the-queue subtasks can either be left out of the reshuffling pool or included.
 - If they are included, provider can be deprived of his work-in-progress.
 - If they are left out, provider can continue work on his work-in-progress across rounds, and be guaranteed a payout, assuming he does finish eventually.
The latter is recommended in general, to account for slow providers/difficult task, however some preemption rule may be necessary.

## Cheap reshuffling challenges

I.e. how to prevent the cost of doing the challenges every round from being to high.
There are two ways:

### Off-chain challenges

In the on-chain challenge scenario, as outlined above, the challenge and the responses are pushed to chain on every reshuffling.
This can lead to (moderate!) expense of gas, compared to hitting the chain only on start and finish.
It also will make every reshuffling last for a few blocks, ca 1-2 minutes, which is some delay, albeit moderate.

There is however a remedy for this.
The requestor and providers `P1`, ..., `Pn` for this task can form a simple ad-hoc P2P network and challenge/publish results off-chain.
If they manage to exchange (signed) commitments required by the protocol this way, they can reshuffle the pool and continue.
If this fail, e.g. challenger doesn't get all commitments in under 2 seconds, he can always decide to publish on chain.

If the off-chain procedure is "chosen", the signed commitments should be possible to be pushed to chain, in case any party misbehaves further down the process.

**NOTE** This remedy is slightly speculative and would be challenging to implement.
There's lots of caveats and opportunities for byzantine behavior.
This needs lots of analysis and careful design.

### Adaptive pool

Other way to further reduce the gas expense for challenges is to adapt the next assignment matrix based on the results from challenged round.
Providers who perform more tasks in round 1 are assigned more tasks in round 2.
This is to prevent the pessimistic scenario when one provider is much stronger than many remaining ones, which would make him get only one equal portion of the pool on each one of numerous reshufflings.

Adaptive pool has some potential to prevent the "Sybil-ing providers" attack, in addition to cutting gas usage.

## Attacks with manipulating the price

### Attack 1: Sybil-ing providers

A possible attack, which is most concerning is this:

1. `P1`, ..., `Pn` are multiple identities of the same entity, e.g. computing farm.
They, and only they (see below for discussion on monopolies) are assigned to a requestor's task.
2. `P1` is doing all the work and publishes challenges, while `P2`, ... `Pn` pretend to be weak laggards and publish nothing

This way the price is pumped, since subtasks of the laggards extend the average subtasks time and final payout.

There are several counter-arguments and potential counter-measures to discuss:

#### Monopolies

The setup where the attacker manages to get only his provider to the task require him to hold monopoly.
If there isn't a monopoly, and the requestor has a choice of reputable providers, then even if a portion of providers signed up perform this attack, the true beneficiaries of the attack are the providers outside attacker's control.
This is because the have the same higher price advantage, but they do not incur costs related to introducing fake laggards and attacking (Paweł's "inverted tragedy of commons" argument).
They will out-speed the attacker and gain much more that him, in terms of final payout.

The assumption of no-monopoly needs to be discussed.

#### Adaptive pool

The no-monopoly effect can be further strengthened by adaptive pool mechanism.
Eventually the majority of subtasks would get assigned to a single provider, diminishing the attack's impact.

#### Adaptive price

Payout for completed subtasks can be gradually reduced on every round.
This provides further motivation to race and weakens the Sybil attack impact.

#### Staking GNT to sign up for task

If a fixed per-address deposit in GNT is required for every provider's address which signs up for a task, the capital cost of doing the Sybil attack is considerable.
When the attacker is forced to hold lots of GNT to do the attack, he would be hurting the network and driving price of GNT down, thus hurting himself.
So there is an incentive to play nice and race for subtasks.

#### Pumping task difficulty drives per hour price down

In the worst case scenario, when this attack is used anyway, an argument can be made, that it will cause the market to push the price per hour down, thereby eliminating the impact on requestors' expenses.
In other words -- the market will discount the risk in the prices, and everyone will be happy anyway, as the computations will be cheaper (for requestors) and more profitable (for providers) than outside Golem.

#### Reputation system

Should also help, naturally.

### Attack 2: exploiting atomic swap

The atomic swap approach would need to be modified to be compatible with the proposed solution.
Otherwise the following attack opens:

1. Requestor controls a subset of providers within the signed-up set `P1`, ..., `Pn`, say only `P1` is honest
2. Requestor will reject all results from **easy** tasks from `P1` -- on grounds of rights the atomic swap gives him -- claiming that results are incorrect
At the same time he will reject (and never spend time on) **difficult** results assigned to his own `P2`, ..., `Pn`
3. This will make the easy tasks be assigned to his providers, which maximizes `P2`, ..., `Pn` payout, robbing `P1` who does the hard work

This attack is possible, because atomic swap opens a possibility of manipulating the task order (requestor's providers disregard difficult tasks). 
(order matters!)
There can't be the possibility of skipping tasks, dropping them from the queue and into the reshuffling pool.

The fix for atomic swap is this: never allow to reject only a single subtask -- if requestor rejects a subtask in atomic swap, he must reject other subtasks from same provider too (and, more rigorously, drop that provider from the list).
This is OK, because incorrect calculation of a subtask is expected to be punished severely.
If the requestor allowed for that provider to continue, he would have acted inconsistently.

That modification must be analyzed, whether or not it opens an attack vector, as the requestor has then a tool to manipulate the providers list **after** the task started.
