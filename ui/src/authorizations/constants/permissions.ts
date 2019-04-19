import {Permission, PermissionResource} from '@influxdata/influx'

// Types

export const ALL_ACCESS_PERMISSIONS = [
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Authorizations},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Authorizations},
  },
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Buckets},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Buckets},
  },
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Dashboards},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Dashboards},
  },
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Sources},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Sources},
  },
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Tasks},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Tasks},
  },
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Telegrafs},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Telegrafs},
  },
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Users},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Users},
  },
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Variables},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Variables},
  },
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Scrapers},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Scrapers},
  },
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Secrets},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Secrets},
  },
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Labels},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Labels},
  },
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Views},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Views},
  },
  {
    action: Permission.ActionEnum.Read,
    resource: {type: PermissionResource.TypeEnum.Documents},
  },
  {
    action: Permission.ActionEnum.Write,
    resource: {type: PermissionResource.TypeEnum.Documents},
  },
]
