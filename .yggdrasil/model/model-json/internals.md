# Model JSON — Internals

## Logic

### ToMotion Segment Parsing
The `Segments` array in `motion3.json` is a flat float64 array where:
1. First two values are the initial point (time, value)
2. Subsequent groups start with a type indicator (0=Linear, 1=Bezier, 2=Stepped, 3=InverseStepped)
3. Type-dependent remaining values:
   - Linear: time, value (3 total: type + 2)
   - Bezier: cp1_time, cp1_val, cp2_time, cp2_val, end_time, end_val (7 total: type + 6)
   - Stepped: end_time, end_value (3 total: type + 2)
   - InverseStepped: end_time, end_value (3 total: type + 2)

### Stepped vs InverseStepped
- **Stepped**: holds the start value until the end time, then jumps
- **InverseStepped**: holds the end value from the start, then changes

## Decisions
- Chose flat array iteration over nested JSON: matches the Cubism SDK spec where segments are encoded as a flat array for compactness
- Chose -1.0 as "use motion fade" sentinel: matches the convention where negative fade times are not physically meaningful
