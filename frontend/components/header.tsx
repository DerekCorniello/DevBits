import React from "react";
import { Header } from "@rneui/themed";
import { StatusBar, Image, StyleSheet } from "react-native";
import { SafeAreaProvider } from "react-native-safe-area-context";

export function MyHeader() {
  return (
    <>
      <StatusBar barStyle="light-content" backgroundColor="red" />
      <Header
        backgroundColor="#151515"
        backgroundImageStyle={{}}
        barStyle="default"
        centerComponent={
          <Image
            source={{
              uri: "https://assets.static-upwork.com/org-logo/950305966420918272",
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
