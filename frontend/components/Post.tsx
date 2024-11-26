import { CheckBox } from "@rneui/themed";
import { Card } from "@rneui/themed";
import { ThemedText } from "@/components/ThemedText";
import { useState } from "react";
import { Stack } from "expo-router";

function Like() {
  const [checked, setState] = useState(true);
  const toggleCheckbox = () => setState(!checked);
  return (
    <CheckBox
      checked={false}
      checkedIcon="octicon DotFillIcon"
      uncheckedIcon="octicon DotIcon"
      checkedColor="#375388"
      onPress={toggleCheckbox}
    />
  );
}

export type PostProps = {
  ID: number;
  User: number; // Linked to User
  Project: number; // Linked to Project
  Likes: number;
  Content: string;
  Comments: number[];
  CreationDate: Date;
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
  return (
    <Card>
      <Card.Title>\{User}</Card.Title>
      <Card.Divider />
      <Card.Image
        style={{ padding: 0 }}
        source={{
          uri: "https://awildgeographer.files.wordpress.com/2015/02/john_muir_glacier.jpg",
        }}
      />
      <Card.Divider />
      <ThemedText type="defaultSemiBold">{Content}</ThemedText>
      <Stack row-align="center">
        <Like />
        <ThemedText type="subtitle">{Likes}</ThemedText>
        <ThemedText type="subtitle">{CreationDate.toUTCString()}</ThemedText>
      </Stack>
    </Card>
  );
}
