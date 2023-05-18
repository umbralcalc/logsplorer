# magicarp

'magicarp' builds on the capabilities and framework provided by the [stochadex package](https://github.com/umbralcalc/stochadex) to enable maximum a posteriori (MAP) optimisation, dynamical process emulation and approximate inference of generalised stochastic models. 

**Notes for future**
- design MAP optimiser algorithm around a local fisher expansion proposal distribution 
- need to think carefully about the framework for generating objective functions so that it can be both general and include known distributions 
- use lots of examples to figure out the best core methods to add in
- include variational inference of the local fisher expansion to get approximate posteriors (could call this Bayesian-ish sensitivity analysis)

## Need more context and documentation?

The design, development, implementation details behind and practical documentation for 'magicarp' can all be found in **Part 2** of this very delightful book: [Diffusing Ideas](https://umbralcalc.github.io/diffusing-ideas/).
