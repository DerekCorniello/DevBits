import React from "react";
import { View, StyleSheet } from "react-native";
import { Header, Icon } from "@rneui/themed";

const TopBar = () => {
  return (
    <View style={styles.container}>
      <Header
        backgroundColor="#ffffff" // Customize color
        leftComponent={
          <Icon
            name="menu"
            type="feather"
            color="#000"
            onPress={() => console.log("Menu clicked!")}
          />
        }
        centerComponent={{ text: "BlueSky", style: styles.title }}
        rightComponent={
          <Icon
            name="search"
            type="feather"
            color="#000"
            onPress={() => console.log("Search clicked!")}
          />
        }
        containerStyle={styles.headerContainer}
      />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#f4f4f4",
  },
  headerContainer: {
    borderBottomWidth: 1,
    borderBottomColor: "#e0e0e0", // Subtle border like BlueSky
  },
  title: {
    color: "#000",
    fontSize: 18,
    fontWeight: "bold",
  },
});

export default TopBar;
