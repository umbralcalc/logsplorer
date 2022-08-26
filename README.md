# mapig

[<img src="_other_resources/diffusing-ideas-badge.svg" width="100">](https://umbralcalc.github.io/)

Maximum a posteriori (MAP) optimisation and approximate inference tools for generalised 2D spatial stochastic models 

**Notes for future**
- design MAP optimiser algorithm around a local fisher expansion proposal distribution 
- need to think carefully about the framework for generating objective functions so that it can be both general and include known distributions 
- use lots of examples to figure out the best core methods to add in
- include variational inference of the local fisher expansion to get approximate posteriors (could call this Bayesian-ish sensitivity analysis)
