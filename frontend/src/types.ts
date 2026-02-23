export interface User {
  id: string
  email: string
  username: string
  created_at: string
  updated_at: string
}

export type EventStatus = 'upcoming' | 'attending' | 'maybe' | 'declined'

export interface Event {
  id: string
  owner_id: string
  title: string
  description: string
  location: string
  starts_at: string
  ends_at: string
  is_all_day: boolean
  color: string
  created_at: string
  updated_at: string
  attendees?: EventAttendee[]
  owner?: User
}

export interface EventAttendee {
  event_id: string
  user_id: string
  status: EventStatus
  is_organizer: boolean
  joined_at: string
  user?: User
}

export type InvitationStatus = 'pending' | 'accepted' | 'declined' | 'expired'

export interface Invitation {
  id: string
  event_id: string
  inviter_id: string
  invitee_email: string
  invitee_id?: string
  token?: string
  status: InvitationStatus
  message: string
  created_at: string
  expires_at: string
  responded_at?: string
  event?: Event
  inviter?: User
}

export interface EventFilter {
  q?: string
  from?: string
  to?: string
  location?: string
  status?: EventStatus
  page?: number
  limit?: number
}

export interface PaginatedEvents {
  events: Event[]
  total: number
  page: number
  limit: number
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
  expires_in: number
  user: User
}

export interface ParsedEvent {
  title: string
  location: string
  starts_at?: string
  ends_at?: string
}

export interface TimeSuggestion {
  starts_at: string
  ends_at: string
  reasoning: string
}
