# DPU/IPU's BMC

Runs redfish server, example:

```text
$ curl --fail http://127.0.0.1:8001/redfish/v1                                                                                                             {
    "@odata.id": "/redfish/v1/",
    "@odata.type": "#ServiceRoot.v1_6_0.ServiceRoot",
    "AccountService": {
        "@odata.id": "/redfish/v1/AccountService"
    },
    "CertificateService": {
        "@odata.id": "/redfish/v1/CertificateService"
    },
    "Chassis": {
        "@odata.id": "/redfish/v1/Chassis"
    },
    "EventService": {
        "@odata.id": "/redfish/v1/EventService"
    },
    "Id": "RootService",
    "Links": {
        "Sessions": {
            "@odata.id": "/redfish/v1/SessionService/Sessions"
        }
    },
    "Managers": {
        "@odata.id": "/redfish/v1/Managers"
    },
    "Name": "Root Service",
    "Oem": {},
    "RedfishVersion": "1.6.0",
    "Registries": {
        "@odata.id": "/redfish/v1/Registries"
    },
    "SessionService": {
        "@odata.id": "/redfish/v1/SessionService"
    },
    "Systems": {
        "@odata.id": "/redfish/v1/Systems"
    },
    "Tasks": {
        "@odata.id": "/redfish/v1/TaskService"
    },
    "UUID": "92384634-2938-2342-8820-489239905423",
    "UpdateService": {
        "@odata.id": "/redfish/v1/UpdateService"
    }
}
```