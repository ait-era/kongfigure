# Kongfigure
Tool to easily configure Kong from JSON resources.

## Supported Kong Resources
 * Services
 * Routes
 * Plugins
 * Consumers
 * Consumer's Credentials
 
The resources are created in the order listed above. This is important if you have resources that depends on each other.

## How to Use
You'll need to create a hierarchy of folders for the resources you want to create.

You need to create a folder for each type of resource you want to setup. The folder name __must__ be in lowercase. Put the JSON file with your resource definition in that folder.
For Consumer's Credentials, in the Consumers folder, create a folder with the name of consumer and then in that folder put the JSON file with the name of the credential.

Example:
```
├── consumers
│   ├── sso_jwt
│   │   └── jwt.json
│   └── sso_jwt.json
├── plugins
│   ├── cors.json
│   ├── jwt.json
│   └── my-custom-plugin.json
├── routes
│   └── root.json
└── services
    └── root.json
```
### How to Create JSON Resource
JSON resources are literally the raw representation of the Kong Resources. Kongfigure use the ID set in the resource to check if it needs to create or update the resource. The only case where it doesn't use the ID directly is for Consumer's Crendeitals where it use the name of the Crendential resource file.


## How to Run Kongfigure
Running Kongfigure is pretty straightforward. You need to pass the kong admin url and the path where all your resources folder are.
`kongfigure --kong-configs dev/kong --kong-url http://kong-admin:8001`
You can use `--dry-run` to see the output of what would created or updated.
```bash
2019/09/30 11:57:08 Kongfigure would patch resource services with id fc9631b1-72f8-4524-a7bf-f85efb65bccd
2019/09/30 11:57:08 Kongfigure would patch resource routes with id e0372fc0-ea10-4141-a130-f0a8555cc8c6
2019/09/30 11:57:08 Kongfigure would patch resource consumers with id 933e9a10-26ba-46d2-9450-075a5a42674f
2019/09/30 11:57:08 Kongfigure would patch resource plugins with id 70433b08-00c9-4031-940a-678b43ce34b6
2019/09/30 11:57:08 Kongfigure would patch resource plugins with id f22d2c5f-2ea5-48d1-9b1b-b5493107697c
2019/09/30 11:57:08 Kongfigure would patch resource plugins with id f3e9b76f-54ba-4535-ab9b-2af6d7513740
2019/09/30 11:57:08 Kongfigure would patch resource consumers/cognito_jwt/jwt with id 285ba56f-36bc-421c-a373-99a21b9cdc76
```