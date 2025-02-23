export interface User {
  id: number
  first_name: string
  last_name: string
  email: string
  is_admin: boolean
  role: Role
  last_login?: string
  created_at: string
  updated_at: string
  deleted_at?: string
}

export interface Role {
  id: number
  name: string
}

export interface Author {
  id: number
  first_name: string
  last_name: string
  slug: string
  bio: string
  user_id: number
  user?: User
  created_at: string
  updated_at: string
  deleted_at?: string
}

export interface Article {
  id: number
  title: string
  slug: string
  short_description: string
  content: string
  status: string
  author_id: number
  author?: Author
  published_at?: string
  created_at: string
  updated_at: string
  deleted_at?: string
}

export interface Image {
  id: number
  url: string
  created_at: string
  deleted_at?: string
}

export interface LoginCredentials {
  email: string
  password: string
}

export interface RegisterData {
  first_name: string
  last_name: string
  email: string
  password: string
}

export interface AuthResponse {
  token: string
  user: User
}
