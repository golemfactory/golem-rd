seqmul = function(start, stop)
    res = 1
    multiplier = start
    while multiplier >= stop
        res *= multiplier
        multiplier -= 1
    end
    res
end

f = function(k, kstar, kT, d)
    res = binomial(big(kstar), big(d))
    res *= seqmul(big(kT), big(kT-d+1))
    res *= seqmul(big(k-kstar), big(k-kstar-kT+d+1))
    res /= seqmul(big(k), big(k-kT+1))
    Float64(res)
end

psybils = function(k, kstar, kT)
    ret = 1 - f(k, kstar, kT, 0) - f(k, kstar, kT, 1)
end

using Distributions

probofatmost = function(top, Tpow, Rpow, Ppow, T, atmost)
    cdf(Binomial(Int32(round(top*Tpow/T/Rpow)), Ppow / Tpow), atmost)
end

probofdissappointment = function(top1, top2, Tpow, Rpow, Ppow, T)
    P1 = probofatmost(top1, Tpow, Rpow, Ppow, T, 0)
    P2 = probofatmost(top2, Tpow, Rpow, Ppow, T, 0.9*top2/T/Rpow*Ppow)
    P1, P2
end

# early adopter
probofdissappointment(12,12*7,100,40,40,10)

# small render farm
probofdissappointment(12,12*7,1000,10,80,10)

# threshold cost of holding deposit to make P-Sybils unprofitable. K*=2
((Ppow, Tpow, C, k, kT) -> Ppow/Tpow * 2 * C * (kT-1)*(k-kT)/(k-1)/(k-2)*(kT-1)/kT^2)(100, 1000, 0.26, 1000, 100)

    # adjust for yearly discount rate 1%, task lasting 10hours * (365*24)/10/0.01
    
    
# relative reduction of cost (as multiples of F) when using R-Sybils. K*=2
((k, kT) -> f(k, 2, kT, 0) + f(k, 2, kT, 1) * ((kT-1)/kT + 1/kT^2) + f(k,2,kT,2)*((kT-2)/kT + 2/kT^2))(100, 10)

# same for K*=1
((k, kT) -> f(k, 1, kT, 0) + f(k, 1, kT, 1) * ((kT-1)/kT + 1/kT^2))
 