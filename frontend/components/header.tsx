import React from "react";
import { Header } from "@rneui/themed";
import { StatusBar, Image, StyleSheet } from "react-native";
import { SafeAreaProvider } from "react-native-safe-area-context";
import { useThemeColor } from "@/hooks/useThemeColor";

export function MyHeader() {
  const cardBackgroundColor = useThemeColor(
    { light: "#C9C9CD", dark: "#151515" },
    "background"
  );
  return (
    <>
      <StatusBar barStyle="light-content" backgroundColor="red" />
      <Header
        backgroundColor={cardBackgroundColor}
        backgroundImageStyle={{}}
        barStyle="default"
        centerComponent={
          <Image
            source={{
              uri: "https://drive.google.com/uc?export=view&id=1KEfYLEUI8mWQtIUGPZXOR5BCfvid49Ul",
            }}
            style={{ width: 200, height: 50 }}
          />
        }
        centerContainerStyle={{}}
        containerStyle={{}}
        leftComponent={{ icon: "menu", color: "grey", size: 35 }}
        leftContainerStyle={{}}
        linearGradientProps={{}}
        rightComponent={{ icon: "terminal", color: "grey", size: 35 }}
        rightContainerStyle={{}}
        statusBarProps={{}}
      />
    </>
  );
}
