export interface File {
  path: string
  size: number
  size_chunks: number
  completed_chunks: number
  frozen_path: string
  priority: number
  is_created: number
  is_open: number
}
