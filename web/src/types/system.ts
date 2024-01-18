export interface System {
  api_version: string
  client_version: string
  library_version: string
  hostname: string
  pid: number
  throttle_global_down_total: number
  throttle_global_up_total: number
  throttle_global_down_rate: number
  throttle_global_up_rate: number
  throttle_global_down_max_rate: number
  throttle_global_up_max_rate: number
}

export interface SystemResponse {
  status: 'ok' | 'error'
  system: System
}

export function IsSystemResponse(object: any): object is SystemResponse {
  return 'system' in object
}
