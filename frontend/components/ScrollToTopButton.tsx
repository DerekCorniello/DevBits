import React from "react";
import { StyleSheet, ScrollView, TouchableOpacity } from "react-native";
import { Icon } from "@rneui/themed";

interface ScrollToTopButtonProps {
  scrollViewRef: React.RefObject<ScrollView>;
}

const ScrollToTopButton: React.FC<ScrollToTopButtonProps> = ({ scrollViewRef }) => {
  const scrollToTop = () => {
    scrollViewRef.current?.scrollTo({ y: 0, animated: true});
    
  };

  return (
    <TouchableOpacity style={styles.button} onPress={scrollToTop} activeOpacity={1}>
      <Icon name="arrow-upward" type="material" color="black" size={25} />
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  button: {
    backgroundColor: "#16ff00",
    padding: 10,
    borderRadius: 30,
    zIndex: 2,
  },
});

export default ScrollToTopButton;
