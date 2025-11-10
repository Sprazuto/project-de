<script setup>
import avatar1 from '@images/avatars/avatar-1.png'
import { useAuth } from '@/composables/useAuth'

// Auth functionality
const auth = useAuth()

// Logout function with debug
const handleLogout = () => {
  // Call logout function
  auth
    .logout()
    .then(() => {})
    .catch((error) => {
      console.error('❌ Logout error from UserProfile:', error)
    })
}
</script>

<template>
  <VBadge dot location="bottom right" offset-x="3" offset-y="3" bordered color="success">
    <VAvatar class="cursor-pointer" color="primary" variant="tonal">
      <VImg :src="avatar1" />

      <!-- SECTION Menu -->
      <VMenu activator="parent" width="230" location="bottom end" offset="14px">
        <VList>
          <!-- 👉 User Avatar & Name -->
          <VListItem>
            <template #prepend>
              <VListItemAction start>
                <VBadge dot location="bottom right" offset-x="3" offset-y="3" color="success">
                  <VAvatar color="primary" variant="tonal">
                    <VImg :src="avatar1" />
                  </VAvatar>
                </VBadge>
              </VListItemAction>
            </template>

            <VListItemTitle class="font-weight-semibold">
              {{ auth.userName }}
            </VListItemTitle>
            <VListItemSubtitle>
              <VChip
                v-if="auth.isAuthenticated.value"
                color="success"
                variant="tonal"
                size="small"
                prepend-icon="tabler-shield-check"
              >
                Authenticated
              </VChip>
              <span v-else>Guest</span>
            </VListItemSubtitle>
          </VListItem>

          <VDivider class="my-2" />

          <!-- 👉 Profile -->
          <VListItem link>
            <template #prepend>
              <VIcon class="me-2" icon="tabler-user" size="22" />
            </template>

            <VListItemTitle>Profile</VListItemTitle>
          </VListItem>

          <!-- 👉 Settings -->
          <VListItem link>
            <template #prepend>
              <VIcon class="me-2" icon="tabler-settings" size="22" />
            </template>

            <VListItemTitle>Settings</VListItemTitle>
          </VListItem>

          <!-- 👉 Dashboard -->
          <VListItem link>
            <template #prepend>
              <VIcon class="me-2" icon="tabler-dashboard" size="22" />
            </template>

            <VListItemTitle>Dashboard</VListItemTitle>
          </VListItem>

          <!-- Divider -->
          <VDivider class="my-2" />

          <!-- 👉 Logout -->
          <VListItem @click="handleLogout">
            <template #prepend>
              <VIcon class="me-2" icon="tabler-logout" size="22" />
            </template>

            <VListItemTitle>Logout</VListItemTitle>
          </VListItem>
        </VList>
      </VMenu>
      <!-- !SECTION -->
    </VAvatar>
  </VBadge>
</template>
