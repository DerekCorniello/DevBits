import { ThemedText } from "@/components/ThemedText";
import TopBar from "@/components/ui/TopBar";
import User from "@/components/User";
import { UserProps, ProjectProps, CommentProps } from "@/constants/Types";
import { View, Text, StyleSheet, Animated } from "react-native";
import React from "react";
import ScrollView from "@/components/ScrollView";
const username = "dev_user1";

export default function ProfileScreen() {
  const [currentUser, setCurrentUser] = React.useState<UserProps | null>(null);
  const [userFollowers, setCurrentUserFollowers] = React.useState<
    number | null
  >(null);
  const [userFollows, setCurrentUserFollows] = React.useState<number | null>(
    null
  );
  // const [currentProjects, setCurrentProjects] =
  // React.useState<ProjectProps | null>(null);

  React.useEffect(() => {
    async function fetchUser() {
      const response = await fetch(`http://localhost:8080/users/${username}`);
      const userData = await response.json();
      console.log(userData);
      setCurrentUser(userData);
    }
    fetchUser();
  }, []);

  React.useEffect(() => {
    async function fetchFollowers() {
      const response = await fetch(
        `http://localhost:8080/users/${username}/followers`
      );
      const userFollowers = await response.json();
      console.log(userFollowers);
      setCurrentUserFollowers(userFollowers);
    }
    fetchFollowers();
  }, []);

  React.useEffect(() => {
    async function fetchFollows() {
      const response = await fetch(
        `http://localhost:8080/users/${username}/follows`
      );
      const userFollows = await response.json();
      console.log(userFollows);
      setCurrentUserFollows(userFollows);
    }
    fetchFollows();
  }, []);

  // React.useEffect(() => {
  //   async function fetchProjects() {
  //     const response = await fetch(
  //       `http://localhost:8080/projects/${username}`,
  //     );
  //     const projectData = await response.json();
  //     console.log(projectData);
  //     setCurrentProjects(projectData);
  //   }
  //   fetchProjects();
  // }, []);

  return (
    <>
      {/* <TopBar /> */}
      {currentUser && (
        <ScrollView>
          <User
            username={currentUser.username}
            bio={currentUser.bio}
            links={currentUser.links}
            created_on={currentUser.created_on}
            picture={currentUser.picture}
          />
          <View>
            <ThemedText type="default">
              <b>{userFollowers ?? 0}</b> followers • <b>{userFollows ?? 0}</b>{" "}
              following
              {/* <b>{projects}</b> streams */}
            </ThemedText>
          </View>
          <ThemedText type="default">Projects</ThemedText>
        </ScrollView>
      )}
    </>
  );
}
