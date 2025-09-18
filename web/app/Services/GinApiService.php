<?php

namespace App\Services;

use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Http;

class GinApiService
{
    private $apiUrl;
    private $accessToken;
    private $refreshToken;

    public function __construct()
    {
        $this->apiUrl = env('GIN_API_URL', 'http://localhost:9000/v1');
    }

    private function getToken()
    {
        $cachedToken = Cache::get('gin_access_token');
        if ($cachedToken) {
            $this->accessToken = $cachedToken;
            return true;
        }

        $email = env('GIN_USER_EMAIL');
        $password = env('GIN_USER_PASSWORD');
        $name = env('GIN_USER_NAME');

        if (!$email || !$password) {
            throw new \Exception('GIN API credentials not set in .env');
        }

        $response = Http::post($this->apiUrl . '/user/login', [
            'email' => $email,
            'password' => $password,
        ]);

        if ($response->successful()) {
            $data = $response->json();
            $this->accessToken = $data['token']['access_token'];
            $this->refreshToken = $data['token']['refresh_token'] ?? null;

            Cache::put('gin_access_token', $this->accessToken, 3600); // 1 hour

            return true;
        } elseif ($response->status() === 401) {
            // User not registered, try register
            $regResponse = Http::post($this->apiUrl . '/user/register', [
                'email' => $email,
                'password' => $password,
                'name' => $name,
            ]);

            if ($regResponse->successful()) {
                // Now login
                $loginResponse = Http::post($this->apiUrl . '/user/login', [
                    'email' => $email,
                    'password' => $password,
                ]);

                if ($loginResponse->successful()) {
                    $data = $loginResponse->json();
                    $this->accessToken = $data['token']['access_token'];
                    $this->refreshToken = $data['token']['refresh_token'] ?? null;

                    Cache::put('gin_access_token', $this->accessToken, 3600);

                    return true;
                }
            }
        }

        throw new \Exception('Failed to authenticate or register with Gin API: ' . $response->body());
    }

    public function getAccessToken()
    {
        if (!$this->accessToken) {
            $this->getToken();
        }
        return $this->accessToken;
    }

    public function call($method, $endpoint, $data = [])
    {
        $this->getToken();

        $response = Http::withHeaders([
            'Authorization' => 'Bearer ' . $this->accessToken,
            'Content-Type' => 'application/json',
        ])->$method($this->apiUrl . $endpoint, $data);

        if ($response->status() === 401) {
            // Token expired, try refresh
            $this->refreshToken();
            $response = Http::withHeaders([
                'Authorization' => 'Bearer ' . $this->accessToken,
                'Content-Type' => 'application/json',
            ])->$method($this->apiUrl . $endpoint, $data);
        }

        return $response;
    }

    private function refreshToken()
    {
        if (!$this->refreshToken) {
            $this->getToken(); // Re-login
            return;
        }

        $response = Http::post($this->apiUrl . '/token/refresh', [
            'refresh_token' => $this->refreshToken,
        ]);

        if ($response->successful()) {
            $data = $response->json();
            $this->accessToken = $data['access_token'];
            $this->refreshToken = $data['refresh_token'];

            Cache::put('gin_access_token', $this->accessToken, 3600);
        } else {
            $this->getToken(); // Re-login if refresh fails
        }
    }

    // Example methods for common endpoints
    public function getArticles()
    {
        return $this->call('get', '/articles');
    }

    public function createArticle($articleData)
    {
        return $this->call('post', '/article', $articleData);
    }
}
