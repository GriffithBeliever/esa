import client from './client'
import type { AuthTokens, User } from '../types'

export const register = (data: { email: string; username: string; password: string }) =>
  client.post<User>('/auth/register', data).then((r) => r.data)

export const login = (data: { email: string; password: string }) =>
  client.post<AuthTokens>('/auth/login', data).then((r) => r.data)

export const logout = () => client.post('/auth/logout')

export const getMe = () => client.get<User>('/auth/me').then((r) => r.data)
