# agents.md

## General Behavior
- Be concise outside the code
- Prefer working code over theory
- Avoid unnecessary abstractions
- Ask clarifying questions only when truly blocked

## Primary Goal
- Produce code that is correct, runnable, and easy for a beginner to follow

## Code Style
- Prefer Go unless another language is explicitly requested
- Write idiomatic code, but favor readability over cleverness
- Keep structure simple and explicit
- Avoid hidden behavior, heavy indirection, and premature abstraction
- Minimize dependencies

## Documentation Requirements
- Overdocument the code
- Assume a new developer may read it with little context
- Add clear comments explaining:
  - what each file is for
  - what each major section does
  - why important decisions were made
  - what inputs and outputs are expected
  - what external libraries are doing
  - what configuration values mean
- For non-obvious code, explain both the purpose and the flow
- Prefer too much explanation over too little
- Use beginner-friendly wording in comments
- Add short step-by-step notes near complicated logic
- Include a short top-of-file overview in main source files

## Output Format
- Default to a minimal runnable example unless a larger structure is requested
- Include all required supporting files such as config files, sample env files, or README notes
- Include build and run instructions
- Keep prose outside the code brief

## Architecture Preferences
- Simple beats complex
- Explicit beats magic
- CLI-first unless UI is requested
- Use config files where appropriate for sessions, connections, and environment-specific settings
- Do not hardcode secrets unless explicitly asked

## Error Handling
- Include basic logging
- Fail fast on critical startup errors
- Print useful runtime output that helps a beginner understand what is happening

## Iteration Style
- First pass should be a complete working baseline
- Refinements can improve structure after functionality is proven

## Nice-to-Have
- Where low cost, include small quality-of-life improvements
- Keep optional enhancements clearly separated from core logic