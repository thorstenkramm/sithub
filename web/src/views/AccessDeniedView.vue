<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="8" lg="6">
        <v-card>
          <v-card-title data-cy="access-denied-title">Access denied</v-card-title>
          <v-card-text>
            Your account is not permitted to access SitHub.
            Please contact your administrator or try a different account.
          </v-card-text>
          <v-card-actions>
            <v-btn data-cy="access-denied-login" href="/oauth/login" variant="tonal">
              Back to login
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
