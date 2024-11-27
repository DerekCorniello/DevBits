import { useState, useRef } from "react";
import { View, Text, Image, StyleSheet, Animated } from "react-native";
import { Card, Icon, CheckBox } from "@rneui/themed";
import { ThemedText } from "@/components/ThemedText";
import { useThemeColor } from "@/hooks/useThemeColor";

function Like() {
  const [checked, setChecked] = useState(false);
  const scaleValue = useRef(new Animated.Value(1)).current;

  const toggleLike = () => {
    Animated.sequence([
      Animated.timing(scaleValue, {
        toValue: 1.1,
        duration: 100,
        useNativeDriver: true,
      }),
      Animated.timing(scaleValue, {
        toValue: 1,
        duration: 200,
        useNativeDriver: true,
      }),
    ]).start();
    setChecked(!checked);
  };
  return (
    <Animated.View style={{ transform: [{ scale: scaleValue }] }}>
      <CheckBox
        center
        containerStyle={[
          styles.checkboxContainer,
          checked && styles.glowEffect,
        ]}
        checkedIcon={
          <Icon
            name="heart-fill"
            type="octicon"
            color="red"
            size={20}
            iconStyle={styles.iconStyle}
          />
        }
        uncheckedIcon={
          <Icon
            name="heart"
            type="octicon"
            color="grey"
            size={20}
            iconStyle={styles.iconStyle}
          />
        }
        checked={checked}
        onPress={toggleLike}
      />
    </Animated.View>
  );
}

export type PostProps = {
  ID: number;
  User: string;
  Project: number;
  Likes: number;
  Content: string;
  CreationDate: Date;
  Comments: number[];
};

export function Post({
  ID,
  User,
  Project,
  Likes,
  Content,
  Comments,
  CreationDate,
}: PostProps) {
  //This changes between light and dark when it is switched. Made it a little off from backround colors.
  const cardBackgroundColor = useThemeColor(
    { light: "light grey", dark: "#000" },
    "background"
  );
  return (
    <Card
      containerStyle={[styles.card, { backgroundColor: cardBackgroundColor }]}
    >
      <View style={styles.header}>
        {/* THis is where  i think the user photo should go*/}
        <ThemedText type="default" style={styles.username}>
          {User}
        </ThemedText>
      </View>
      <ThemedText type="default" style={styles.content}>
        {Content}
      </ThemedText>
      {/*looks weird i dont know how to get rid of the space at the top*/}
      <Card.Divider style={styles.divider} />
      <View style={styles.footer}>
        <Like />
        <Text style={styles.likes}>{Likes} Likes</Text>
      </View>
      <Text style={styles.date}>{CreationDate.toUTCString()}</Text>
    </Card>
  );
}

// I fr dont have the mental capacity to figure out how to make the text colors match the theme so gl
const styles = StyleSheet.create({
  card: {
    borderRadius: 12,
    padding: 10,
    marginBottom: 20,
  },
  header: {
    flexDirection: "row",
    alignItems: "center",
    marginBottom: 10,
  },
  divider: {
    marginTop: 5,
    marginBottom: 5,
  },
  username: {
    fontSize: 20,
    fontWeight: "bold",
  },
  content: {
    fontSize: 14,
    lineHeight: 20,
    marginBottom: 10,
  },
  iconStyle: {
    backgroundColor: "transparent",
  },
  checkboxContainer: {
    padding: 1,
    backgroundColor: "transparent",
    borderWidth: 0,
  },
  glowEffect: {
    shadowColor: "red",
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.8,
    shadowRadius: 10,
  },
  footer: {
    flexDirection: "row",
    alignItems: "center",
    marginTop: 5,
  },
  likes: {
    fontSize: 14,
    color: "grey",
  },
  date: {
    fontSize: 12,
    color: "grey",
    alignSelf: "flex-end",
  },
});
