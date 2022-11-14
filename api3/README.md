# APIClarity APIs 

APIClarity is composed by a set of core functionalities, implemented by the so-called `APIClarity Core`, and by a set of additional functionalities implemented in modules.

APIClarity core exposes a core API specified in [/api3/core/openapi.yaml](/api3/core/openapi.yaml).

Modules backend code is placed under [/backend/pkg/modules/internal/](/backend/pkg/modules/internal/), where each module can have its own folder containing all needed code.

Each module can also expose its own API, which must be specified in `backend/pkg/modules/internal/[MODULE_NAME]/restapi/openapi.yaml`. 
The operations served by modules' APIs are served under the following subpath of APIClarity API `modules/{MODULE_NAME}`. 

For example if API Clarity Core API is served under 
```
https://apiclarity-server/api/
```
then the module "example" will have its APIs exposed under 
```
https://apiclarity-server/api/modules/example/
```

## APIClarity common specification
Both core and modules OAPI specifications may refer to common schemas and types defined in [/api3/core/openapi.yaml](/api3/core/openapi.yaml)
If modules generate code from the specifications then they should use the common package that holds the generated code for the common specifications.
Whenever a change is made to the common specification, the code needs to be generated as follows:
```
cd api3/common; go generate
```
or by using the Makefile as described below.

## APIClarity Global specifications
Core and Modules specifications are aggregated in a single specification to allow for APIClarity clients to have a single API definitions. This is what we call `APIClarity Global API Specification`, which is automatically generated and kept in [/api3/global/openapi.gen.yaml](/api3/global/openapi.gen.yaml) 

The automatic generation of such global specification must be done whenever changes are made to the common specification, to the core specification, or to any module specification. This is done by running the following tool:
```
cd tools/spec-aggregator; go run main.go
```
or by using the Makefile as described below.

Client and server code is also generated for such global spec as follows:
```
cd api3/global; go generate
```
or by using the Makefile as described below.


## APIClarity Notification specification

APIClarity is also able to send notifications to interested registered listeners.
The API that a listener needs to implement is specified in [/api3/notifications/openapi.gen.yaml](/api3/notifications/openapi.gen.yaml), this is automatically generated by aggregating all notifications defined in the core and in the modules.

The automatic generation of such global specification must be done whenever changes are made to the common specification, to the core specification, or to any module specification. This is done by running the following tool:
```
cd tools/spec-aggregator; go run main.go
```
or by using the Makefile as described below.
Client and server code is also generated for such global spec as follows:
```
cd api3/notifications; go generate
```
or by using the Makefile as described below.

## Notification Schemas and conventions

In the notifications specification, all API clarity notifications are part of a single parent object schema definition, called `APIClarityNotification`, where each specific instance is a member of the `OneOf` property of such parent schema.

For example, when this documentation is being written, the automatically generated definition of the parent schema is:
```
    APIClarityNotification:
      discriminator:
        mapping:
          ApiFindingsNotification: '#/components/schemas/ApiFindingsNotification'
          AuthorizationModelNotification: '#/components/schemas/AuthorizationModelNotification'
          TestProgressNotification: '#/components/schemas/TestProgressNotification'
          TestReportNotification: '#/components/schemas/TestReportNotification'
        propertyName: notificationType
      oneOf:
      - $ref: '#/components/schemas/ApiFindingsNotification'
      - $ref: '#/components/schemas/AuthorizationModelNotification'
      - $ref: '#/components/schemas/TestReportNotification'
      - $ref: '#/components/schemas/TestProgressNotification'
```
Note that each of the Notification Definition must include the discriminator field `notificationType`. To achieve this it is mandatory for all specific notifications to be in the form:
```
TestReportNotification:
    allOf: 
    - $ref: '../../../../../../api3/common/openapi.yaml#/components/schemas/BaseNotification'
    - $ref: '#/components/schemas/ShortTestReport'
```
which means that they should have a allOf proprty where the first element of te list corresponds to the BaseNotification type defined in the common specification.

Note also in the mapping section, that by convention the discriminator value for each notification corresponds to the name of the schema definition of that notification, e.g. The discrminator `notificationType` for the notification `#/components/schemas/TestReportNotification` is `TestReportNotification`

These conventions are used by the spec aggregation tool to correctly find the core and module notifications and aggregate them. Failure to follow such conventions will cause a failure in the aggregation.

## Use of Makefile
All spec aggregation and code generation can be executed through the Makefile at the repo root:
```
make api3
```