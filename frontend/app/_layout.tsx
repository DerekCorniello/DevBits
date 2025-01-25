import {
  DarkTheme,
  DefaultTheme,
  ThemeProvider,
} from "@react-navigation/native";
import { useFonts } from "expo-font";
import { Stack } from "expo-router";
import * as SplashScreen from "expo-splash-screen";
import { StatusBar } from "expo-status-bar";
import { useEffect } from "react";
import "react-native-reanimated";
import { SafeAreaView, Platform } from "react-native";
import { useColorScheme } from "@/hooks/useColorScheme";
import { useAuth0, Auth0Provider } from "react-native-auth0";
import { ThemedView } from "@/components/ThemedView";
// Prevent the splash screen from auto-hiding before asset loading is complete
SplashScreen.preventAutoHideAsync();

export default function RootLayout() {
  const colorScheme = useColorScheme();
  const [loaded] = useFonts({
    SpaceMono: require("../assets/fonts/SpaceMono-Regular.ttf"),
  });

  useEffect(() => {
    if (loaded) {
      // Hide splash screen
      SplashScreen.hideAsync();
    }
  }, [loaded]);

  // Render null if fonts are not loaded
  if (!loaded) {
    return null;
  }

  return (
    <Auth0Provider
      domain={"dev-mcbbo7b2hkpnb65f.us.auth0.com"}
      clientId={"2mZsj7PZmvraWwfbtUMYdHndLgqjufK5"}
    >
      <ThemeProvider value={colorScheme === "dark" ? DarkTheme : DefaultTheme}>
        <Stack screenOptions={{ headerShown: false }}>
          <Stack.Screen name="(tabs)" />
          <Stack.Screen name="+not-found" />
        </Stack>
        <StatusBar style={colorScheme === "dark" ? "light" : "dark"} />
      </ThemeProvider>
    </Auth0Provider>
  );
}
