// Package pubro provides a singleton provider for registering publishers
// which allows initialization of those publishers quickly.
/*

Publishers are registered by providing a PubMeta struct:

  pubro.Register({
    Name: "builders/flint",
    Desc: `Flint provides a publisher that ....`
    Inject: func(config FlintConfig) pub.Publisher {
      return ....
    },
  })


Using a registered Publisher from the registry:

  flint := pubro.New("builders/flint",&FlintConfig{})

Including the basics above, pubro comes with a building systems which takes
a registry of publishers and connections meta and produces the appropriate
build order to create the necessary pub structure
*/
package pubro
