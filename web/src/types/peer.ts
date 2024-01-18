export interface Peer {
  peer_id: string
  address: string
  port: number
  client_version: string
  peer_rate: number
  peer_total: number
  up_rate: number
  up_total: number
}
