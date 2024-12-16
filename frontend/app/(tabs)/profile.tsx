import TopBar from "@/components/ui/TopBar";
import User from "@/components/User";
import { UserProps } from "@/constants/Types";
import React from "react";

const username = "dev_user1";

export default function ProfileScreen() {
  const [currentUser, setCurrentUser] = React.useState<UserProps | null>(null);

  React.useEffect(() => {
    async function fetchUser() {
      const response = await fetch(`http://localhost:8080/users/${username}`);
      const userData = await response.json();
      console.log(userData);
      setCurrentUser(userData);
    }
    fetchUser();
  }, []);

  return (
    <>
      {/* <TopBar /> */}
      {currentUser && (
        <User
          username={currentUser.username}
          bio={currentUser.bio}
          links={currentUser.links}
          created_on={currentUser.created_on}
          picture={currentUser.picture}
        />
      )}
    </>
  );
}
