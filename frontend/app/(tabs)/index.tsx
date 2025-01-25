import React, { useState, useRef } from "react";
import { Animated, StyleSheet, ScrollView, SafeAreaView, View } from "react-native";
import { Post } from "@/components/Post";
import CreatePost from "@/components/CreatePost";
import { MyFilter } from "@/components/filter";
import { MyHeader } from "@/components/header";
import ScrollToTopButton from "@/components/ScrollToTopButton";

export default function HomeScreen() {
  const [scrollY] = useState(new Animated.Value(0));
  const scrollViewRef = useRef<ScrollView>(null);
  const headerOpacity = scrollY.interpolate({
    inputRange: [0, 200],
    outputRange: [1, 0],
    extrapolate: "clamp",
  });
  const topscrollOpacity = scrollY.interpolate({
    inputRange: [300, 500],
    outputRange: [0, 1], 
    extrapolate: "clamp",
  });
  return (
    <SafeAreaView style={{ flex: 1 }}>
      <Animated.View style={[styles.header, { opacity: headerOpacity }]}>
        <MyHeader />
      </Animated.View>
      <Animated.View style={[styles.filterContainer, { opacity: headerOpacity }]}>
        <MyFilter />
      </Animated.View>
      <Animated.View style={[styles.scrolltotopbutton, { opacity: topscrollOpacity }]}>
        <ScrollToTopButton scrollViewRef={scrollViewRef} />
      </Animated.View>
      <ScrollView
        ref={scrollViewRef}
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
          likes={Math.floor(Math.random() * 100)}
          content={`This is a test post with a content that I makde for the test post that has content. This content is for testing only. that I made. -E`}
          created_on={`2021-01-01T00:00:00Z`}
          comments={[20]}
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
  scrolltotopbutton: {
    alignSelf: "flex-start",
    position:"absolute",
    paddingTop:60,
    zIndex:2,
    paddingLeft:20
  },
});
