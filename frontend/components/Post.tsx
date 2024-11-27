import { useState, useRef } from "react";
import { View, Text, StyleSheet, Animated } from "react-native";
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

function Comment() {
  return (
    <Icon
      name="comment"
      type="octicon"
      color="grey"
      size={20}
      iconStyle={styles.iconStyle}
    />
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
  const cardBackgroundColor = useThemeColor(
    { light: "light grey", dark: "#151515" },
    "background"
  );
  return (
    <Card
      containerStyle={[styles.card, { backgroundColor: cardBackgroundColor }]}
    >
      <View style={styles.header}>
        <ThemedText type="default" style={styles.username}>
          {User}
        </ThemedText>
        <ThemedText type="subtitle" style={styles.project}>
          Stream {Project}
        </ThemedText>
      </View>
      <ThemedText type="default" style={styles.content}>
        {Content}
      </ThemedText>
      <Card.Divider style={styles.divider} />
      <View style={styles.footer}>
        <View style={styles.actionContainer}>
          <Like />
          <Text style={styles.bottomText}>{Likes} likes</Text>
        </View>
        <View style={styles.actionContainer}>
          <Comment />
          <Text style={styles.bottomText}>{Comments.length} bits</Text>
        </View>
      </View>
      <Text style={styles.date}>
        {CreationDate.toLocaleString("en-US", {
          weekday: "long",
          year: "numeric",
          month: "short",
          day: "numeric",
          hour: "numeric",
          minute: "numeric",
          hour12: true,
        })}
      </Text>
    </Card>
  );
}

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
  project: {
    fontSize: 14,
    marginLeft: 5,
  },
  content: {
    fontSize: 14,
    lineHeight: 20,
    marginBottom: 10,
  },
  iconStyle: {
    backgroundColor: "transparent",
    marginRight: 5,
  },
  checkboxContainer: {
    padding: 1,
    backgroundColor: "transparent",
    borderWidth: 0,
    marginRight: 0,
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
    justifyContent: "space-between",
    marginTop: 5,
  },
  actionContainer: {
    flexDirection: "row",
    alignItems: "center",
  },
  bottomText: {
    fontSize: 14,
    color: "grey",
    marginLeft: 5,
  },
  date: {
    fontSize: 12,
    color: "grey",
    alignSelf: "flex-end",
    marginTop: 5,
  },
});
