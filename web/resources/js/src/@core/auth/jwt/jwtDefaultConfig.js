export default {
    // Endpoints
    loginEndpoint: '/user/login',
    registerEndpoint: '/user/register',
    refreshEndpoint: '/token/refresh',
    logoutEndpoint: '/user/logout',

    // This will be prefixed in authorization header with token
    // e.g. Authorization: Bearer <token>
    tokenType: 'Bearer',

    // Value of this property will be used as key to store JWT token in storage
    storageTokenKeyName: 'access_token',
    storageRefreshTokenKeyName: 'refresh_token',
}
