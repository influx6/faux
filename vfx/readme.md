# VFX
 Vfx is a idiomatic web animation library which brings the style of [VelocityJS](julian.com/research/velocity/)
 to Go. It brings a different approach in the way animation are built but
 has a solid foundation to allow the best performance ensuring to reduce layout
 trashing as much as possible.

## Install

  ```go

    go get -u github.com/influx6/faux/tree/master/vfx

  ```

## Concept

  - Sequence and Writers
  `Sequence` in VFX define a behaviour that changes respectively on every tick of
  the animation clock, they hold the calculations that are necessary to achieve
  the desired behaviour, in a way they define a means of doing a deferred processing
  of the expected behaviours.
  A `sequence` can be a width or height transition, or an opacity animator that
  produces for every iteration the desired end result. Sequences can be of any
  type as defined by the animation creator, which provides a powerful but
  flexible system because multiple sequences can be bundled together to produce
  new ones.
  `Sequences` return `Writers` which are the calculated end result, the reason there
  exists such a concept of Writers is due to the need in reducing the effects of
  massive layout trashing, which is simple updating and reading changing properties
  that is causing recalculation of the display by the browser, which costs in
  performance in rendering. When writers are returned for a sequence or sets of
  sequence then this allows effective batching of these changes which reduces to
  minimum the update calculation performed by the browser DOM.

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
