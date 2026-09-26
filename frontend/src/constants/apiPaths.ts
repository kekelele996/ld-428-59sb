export const apiPaths = {
  artworks: '/api/v1/artworks',
  exhibitions: '/api/v1/exhibitions',
  artists: '/api/v1/artists',
  interactions: '/api/v1/interactions',
  authLogin: '/api/v1/auth/login',
  sessions: (exhibitionId: string) => `/api/v1/exhibitions/${exhibitionId}/sessions`,
  reservations: '/api/v1/reservations',
  exhibitionReservations: (exhibitionId: string) => `/api/v1/exhibitions/${exhibitionId}/reservations`,
  cancelReservation: (reservationId: string) => `/api/v1/reservations/${reservationId}/cancel`,
};
