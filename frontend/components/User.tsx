import ScrollView from "@/components/ScrollView";
import TopBar from "@/components/ui/TopBar";
import { View, Text, StyleSheet, Animated } from "react-native";
import { ThemedText } from "@/components/ThemedText";
import { UserProps } from "@/constants/Types";
import { ExternalPathString, Link } from "expo-router";

export default function User({
  username,
  bio,
  links,
  created_on,
  picture,
}: UserProps) {
  let CreationDate = new Date(created_on);
  return (
    <>
      {/* <TopBar /> */}
      <ScrollView>
        <View style={styles.header}>
          <ThemedText type="default" style={styles.username}>
            {username}
          </ThemedText>
        </View>
        <ThemedText type="default" style={styles.bio}>
          {bio}
        </ThemedText>
        {links.map((link, index) => (
          <Link key={index} href={link as ExternalPathString}>
            <ThemedText type="default" style={styles.link}>
              {link}
            </ThemedText>
          </Link>
        ))}
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
      </ScrollView>
    </>
  );
}

const styles = StyleSheet.create({
  header: {
    flexDirection: "row",
    alignItems: "center",
    marginBottom: 10,
  },
  username: {
    fontSize: 20,
    fontWeight: "bold",
  },
  bio: {
    fontSize: 14,
    lineHeight: 20,
    marginBottom: 10,
  },
  date: {
    fontSize: 12,
    color: "grey",
    alignSelf: "flex-end",
    marginTop: 5,
  },
  link: {
    fontSize: 14,
    color: "#007AFF",
    lineHeight: 20,
    marginBottom: 10,
  },
});
