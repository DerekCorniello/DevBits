import { ThemedText } from "@/components/ThemedText";
import TopBar from "@/components/ui/TopBar";
import User from "@/components/User";
import { UserProps, ProjectProps, CommentProps } from "@/constants/Types";
import React from "react";
import ScrollView from "@/components/ScrollView";
const username = "dev_user1";

export default function ProfileScreen() {
  const [currentUser, setCurrentUser] = React.useState<UserProps | null>(null);
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
          <ThemedText type="default">Projects</ThemedText>
        </ScrollView>
      )}
    </>
  );
}
