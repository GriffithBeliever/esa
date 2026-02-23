import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Layout } from './components/Layout'
import { ProtectedRoute } from './auth/ProtectedRoute'
import { LoginPage } from './auth/LoginPage'
import { RegisterPage } from './auth/RegisterPage'
import { EventsPage } from './events/EventsPage'
import { EventDetail } from './events/EventDetail'
import { CreateEventPage } from './events/CreateEventPage'
import { EditEventPage } from './events/EditEventPage'
import { InvitationsPage } from './invitations/InvitationsPage'
import { InviteAcceptPage } from './invitations/InviteAcceptPage'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 30_000,
    },
  },
})

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          {/* Public routes */}
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/invite/:token" element={<InviteAcceptPage />} />

          {/* Protected routes */}
          <Route element={<ProtectedRoute />}>
            <Route
              path="/*"
              element={
                <Layout>
                  <Routes>
                    <Route index element={<Navigate to="/events" replace />} />
                    <Route path="events" element={<EventsPage />} />
                    <Route path="events/new" element={<CreateEventPage />} />
                    <Route path="events/:id" element={<EventDetail />} />
                    <Route path="events/:id/edit" element={<EditEventPage />} />
                    <Route path="invitations" element={<InvitationsPage />} />
                  </Routes>
                </Layout>
              }
            />
          </Route>

          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
