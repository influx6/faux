# VFX
 Vfx is a idiomatic web animation library which brings the style of [VelocityJS](https://julian.com/research/velocity/)
 to Go. It brings a different approach in the way animations are built but
 has a solid foundation to allow the best performance ensuring to reduce layout
 trashing as much as possible.

## Install

  ```go

    go get -u github.com/influx6/faux/vfx

  ```

## Concept

  - Sequence and Writers

  `Sequence` in VFX define a behaviour that changes respectively on every tick of
  the animation clock, they hold the calculations that are necessary to achieve
  the desired behaviour. In a way, they define a means of providing a
  processing of the deferred expected behaviour.

  A `sequence` can be a width or height transition, or an opacity animator that
  produces for every iteration the desired end result.
  `Sequences` can be of any type as defined by the animation creator, which
  provides a powerful but flexible system because multiple sequences can be bundled
  together to produce new ones.

  `Sequences` return `Writers`, which are the calculated end result for a animation step.
  The reason there exists such a concept, is due to the need in reducing the effects of
  massive layout trashing, which is basically the browser re-rendering of the DOM
  due to massive multiple changes of properties of different elements, which create
  high costs in performance.

  Writers are returned from sequences to allow effective batching of these changes
  which reduces and minimizes the update calculation performed by the browser DOM.

  - Animation Frames

  `Animation frames` in VFX are a collection of sequences which according to a
  supplied stat will produce the total necessary sequence writers need to
  achieve the desired animation within a specific frame or time of the animation
  loop. It is the central organizational structure in VFX.

  - Stats

  Stats in VFX are captions of current measurements of the animation loop and the
  properties for `Animation Frames`, using stats VFX calls all sequence to produce
  their writers by using the properties of the stats to produce the necessary change
  and easing behaviours that is desired to be achieved.

## Usage
