# Deterministic computation
## Why we'd like to have determinism
Deterministic computations are very convenient.
* we can do verification by replication
  *  we don't have to write custom metric per application
* we can [verify-before-swap](https://github.com/imapp-pl/golem_rd/wiki/Verification-of-computation#verification-before-result-delivery)
* results of computation which greatly reduces amount of needed trust between parties

## Problems others have with non-determinism
A holy grail of most distributions are deterministic builds. Debian itself has a debhelper script `dh_strip_nondeterminism`, which gets rid of elements such as dates embedded in files. It provides handlers for normalizing many kinds of binary files, such as `png`.

Chromium discusses deterministic builds [here](https://www.chromium.org/developers/testing/isolated-testing/deterministic-builds)

While deterministic builds are not exactly, what troubles us, these problems give
an important context.

## Sources of non-determinism
All discussion here will tackle mostly Linux environment.
### System-provided randomness
At first look, it seems that it's enough to replace the `/dev/{u,}random` files
with something created by ourselves. This is not the case, though.

The Linux kernel, starting from version 3.19, provides a system call `getrandom`.
This source of randomness is commonly used by modern programming languages such
as Python or Rust.

These languages exploit the syscall in two different ways. CPython, starting from
3.6 uses the glibc `getrandom()` C wrapper whenever possible. The official Rust
crate `rand`, on the other hand, makes direct use of the C `syscall` function,
effectively bypassing all attempts of replacing the wrappers.
### Time
One of the common ways to seed PRNGs is to use time. That's in fact the most
common method in C or Fortran.

There are multiple system calls which provide
getting the current time: `clock_gettime`, `gettimeofday`, `time`. Some of them,
e.g. time are implemented as a part of vDSO, which means they can't be easily
caught using `ptrace`.
### Threads interleaving
One source of non-determinism is the interleave of threads. First idea that comes
to one's mind when considering this problem is serializing the threads. This is
unfortunately not so easy as it appears at the first sight.

One could say that any nondeterminism is a sign of a buggy program but it's not
the case.

Even if we use different synchronization techniques, on Linux mutexes/semaphores
do not guarantee any order the processes being released from barriers. This means,
that many correct programs will return nondeterministic results from sheer
usage of threads.

An example could be a program whose all threads insert elements to a vector `v`
guarded with a mutex. Then the order of the elements will be non-deterministic,
of course, but `v.sort()` will be deterministic.

Even if we decide to try taming the multi threaded programs, the following
considerations have to be taken into account:
* all data access should be serialized in a predetermined manner

### Network connectivity
Even if we patch any other sources of non-determinism, any network-based
interaction may be a source of non-determinism. Even ICMP packets.

Although applications, which currently run on Golem, require no network
connectivity, it may be very important for other types of applications.
Golem project is currently aspiring to a *distributed supercomputer* and one of
the tasks commonly computed on such machines is weather forecast. It is an
extremely computation-intensive task which cannot be performed without
inter-node communication.

Sidenote: inter-node computing in heterogenous environments involves different
kinds of problems, which have nothing to do with security and our out of
the scope of this document.
### File access
Though the Docker environment is fairly standarized, every file access causes
non-determinism. One can read the file timestamps, contents of `/tmp` may change
(and any files created with `mktemp` will have non-deterministic names).
### Recompilation
Mainly problem for deterministic compilation - the ANSI C macros: `__FILE__` and
`__LINE__` may introduce nondeterminism.
### Environmental variables
Many applications are controlled via environmental variables, for example those
using OpenMP. While it's impossible to completely ban using environmental
variables (this would break `execvp`), applications can gain some entropy, for
exampling by
```
import subprocess
subprocess.call(["bash", "-c", "echo $SECONDS"])
```
### System parameters
Fortran MPI-based applications run on supercomputers often suffer from the fact
that the language lacks proper RNG initialization. A common technique is to use
as the seed:
```
PID of the current process * MPI process rank * current time
```

the list follows: device hostname, local IP address, `uname -a`, library versions,
uptime, scheduler parameters...

While many of these issues are tackled by Docker containerization, the following
depend on the host, for example - surprisingly - uptime.

### Undefined behavior
#### Use of uninitialized memory.

Although [kernel zeroes memory](https://stackoverflow.com/questions/6004816/kernel-zeroes-memory) as security measure, it does not really help because you can't really know if you
will receive a new zero'ed page or the old page you've been already using.
Which amounts to one bit of entropy. Other possible source of entropy
is reordering of reused, non-zeroed pages.

#### Arithmetic caveats
Some arithmetic operations, which invoke UB, will behave differently on different
architectures. For example `INT_MIN / -1` will cause a `SIGFPE` on x86 but will
yield a value on ARM.

### ... TODO

On the other hand, the macro `offsetof` defined in ANSI C often invokes undefined
behavior itself without any bad side effects.

## Methods of enforcing determinism

### The `LD_PRELOAD` trick
A simple way of injecting arbitrary code into programs is using `LD_PRELOAD`.
For example:
```
LD_PRELOAD=libJestemSuperHakierem.so ./executable
```
This works, because linker will check if the symbol is available in any of the
preloaded libraries before falling back to the real ones.

This approach can override only the dynamically linked symbols, so is not
suitable for statically-linked applications or those that use direct syscalls.

An example of such hack is presented in folder `ld_preload` (TODO)
### `ptrace`
Since a large part of non-determinism originates from the system calls, a good
approach is to hack the syscalls themselves. This can be done with help of ptrace.

An example of using `ptrace` for mocking randomness is described later.

It should be noted that not all system calls can be mocked using this method.
Some of the system calls are implemented using vDSO. This mean they don't trigger
processor mode switch and `ptrace` won't notify about them.

### Injecting shared libraries to child processes
TODO `dlopen`/`libloading`
## Proof of concept for mocking randomness

@marmistrz wrote a program, `randmockery` which proves that we can use `ptrace`
to mock randomness provided by system calls. It's open source and available
on [GitHub](https://github.com/marmistrz/randmockery). Known issues are
mentioned in the issue tracker.

## Other notes

**More about floating point non determinism:**

Excellent summary from gamedev industry - [Is it possible to make floating point calculations completely deterministic?](https://gafferongames.com/post/floating_point_determinism/)

From BOINC:

https://boinc.berkeley.edu/trac/wiki/HomogeneousRedundancy
https://boinc.berkeley.edu/trac/wiki/ValidationIntro

For scientist's who really wants to dig into- [What Every Computer Scientist Should Know About Floating-Point Arithmetic](http://docs.oracle.com/cd/E19957-01/806-3568/ncg_goldberg.html)
