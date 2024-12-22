import React, { useState } from "react";
import { Animated, StyleSheet, ScrollView, View, SafeAreaView } from "react-native";
import { Post } from "@/components/Post";
import { useThemeColor } from "@/hooks/useThemeColor";
import CreatePost from "@/components/CreatePost";
import { MyFilter } from "@/components/filter";
import { MyHeader } from "@/components/header";

export default function HomeScreen() {
  const [scrollY] = useState(new Animated.Value(0));
  const headerOpacity = scrollY.interpolate({
    inputRange: [0, 200],
    outputRange: [1, 0],
    extrapolate: "clamp",
  });
  const cardBackgroundColor = useThemeColor(
    { light: "#fff", dark: "#040607" },
    "background"
  );
  return (
    <SafeAreaView style={{ flex: 1 }}>
      <Animated.View style={[styles.header, { opacity: headerOpacity }]}>
        <MyHeader />
      </Animated.View>
      <Animated.View style={[styles.filterContainer, {opacity: headerOpacity}]}>
        <MyFilter />
      </Animated.View>
      <ScrollView
        contentContainerStyle={styles.scrollContainer}
        onScroll={Animated.event(
          [{ nativeEvent: { contentOffset: { y: scrollY } } }],
          { useNativeDriver: false }
        )}
        scrollEventThrottle={16}
      >
        <Post
          id={1}
          user={2}
          project={1}
          likes={69}
          content="This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app."
          created_on="2021-01-01T00:00:00Z"
          comments={[]}
        />
                <Post
          id={1}
          user={2}
          project={1}
          likes={69}
          content="This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app."
          created_on="2021-01-01T00:00:00Z"
          comments={[]}
        />
                <Post
          id={1}
          user={2}
          project={1}
          likes={69}
          content="This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app."
          created_on="2021-01-01T00:00:00Z"
          comments={[]}
        />
                <Post
          id={1}
          user={2}
          project={1}
          likes={69}
          content="This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app. This is a test post. It should be displayed in the app."
          created_on="2021-01-01T00:00:00Z"
          comments={[]}
        />
      </ScrollView>
      <CreatePost />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  header: {
    position: "absolute",
    zIndex: 1,
    width: "100%",
  },
  filterContainer: {
    position: "absolute",
    alignSelf: "flex-end",
    padding: 20,
    paddingTop: 140,
    zIndex: 1,
  },
  scrollContainer: {
    paddingTop: 150,
  },
});
