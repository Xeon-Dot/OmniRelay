<template>
  <div id="app">
    <router-view />
    <!-- Single global snackbar host; any component can show one via the store. -->
    <OrSnackbar />
  </div>
</template>

<style scoped>
/*
 * App.vue intentionally carries no global styles.
 * The old version here held a full second copy of the design system (legacy
 * --clr-* variables plus ~200 lines of Vuetify `!important` overrides) which
 * would fight the M3 token layer. Reset, typography, focus and the
 * interactive state layers now live in styles/{base,ui}.css, and component
 * surfaces come from the token files.
 */
#app {
  display: flex;
  min-block-size: 100dvh;
}
</style>
