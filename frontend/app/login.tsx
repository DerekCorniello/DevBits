import React from "react";
import { View, Text, Button, StyleSheet } from "react-native";
import LoginButton from "@/components/LoginButton";
import LogoutButton from "@/components/LogoutButton";
import { useAuth0 } from "react-native-auth0";
import { ThemedView } from "@/components/ThemedView";

export default function LoginScreen() {
  const { user, error } = useAuth0();

  return (
    <ThemedView style={styles.container}>
      <Text style={styles.title}>Welcome to DevBits!</Text>
      {user ? (
        <>
          <Text style={styles.message}>Hello, {user.name}!</Text>
          <LogoutButton />
        </>
      ) : (
        <>
          <Text style={styles.message}>Please log in to continue.</Text>
          <LoginButton />
        </>
      )}
      {error && <Text style={styles.error}>Error: {error.message}</Text>}
    </ThemedView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    padding: 16,
    backgroundColor: "#fff",
  },
  title: {
    fontSize: 24,
    fontWeight: "bold",
    marginBottom: 16,
  },
  message: {
    fontSize: 16,
    marginBottom: 8,
  },
  error: {
    color: "red",
    marginTop: 8,
  },
});
