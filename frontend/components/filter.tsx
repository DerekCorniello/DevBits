import React, { useState } from "react";
import { View, Text, Modal, Animated, Button } from "react-native";
import { FAB, Icon, Switch } from "@rneui/themed";

export const MyFilter: React.FC = () => {
  const [modalVisible, setModalVisible] = useState(false);
  const [slideAnim] = useState(new Animated.Value(-400));
  const [value, setValue] = React.useState(false);

  const toggleModal = () => {
    setModalVisible(!modalVisible);
    if (!modalVisible) {
      slideAnim.setValue(-400);
      Animated.timing(slideAnim, {
        toValue: 0,
        duration: 100,
        useNativeDriver: true,
      }).start();
    } else {
      Animated.timing(slideAnim, {
        toValue: 400,
        duration: 100,
        useNativeDriver: true,
      }).start();
    }
  };
  return (
    <View style={{ backgroundColor: "transparent" }}>
      <FAB
        visible={true}
        title="Filter"
      
        icon={<Icon name="tune" type="material" color="black" size={25} />}
        color="#16ff00"
        size="large"
        onPress={toggleModal}
        titleStyle={{ color: "black",fontSize:15 }}
      />
      <Modal
        animationType="fade"
        transparent={true}
        visible={modalVisible}
        onRequestClose={toggleModal}
      >
        <View
          style={{
            flex: 1,
            justifyContent: "center",
            alignItems: "center",
            backgroundColor: "rgba(32, 48, 32, 0.5)",
          }}
        >
          <Animated.View
            style={{
              width: 300,
              height: 200,
              backgroundColor: "white",
              borderRadius: 10,
              justifyContent: "center",
              alignItems: "center",
              transform: [{ translateY: slideAnim }],
              padding: 20,
            }}
          >
            <Text style={{ fontSize: 20, marginBottom: 20 }}>
              Da Filter Options
            </Text>
            <Switch
              color="#2089dc"
              value={value}
              onValueChange={() => setValue(!value)}
            />
            <Button title="Close Filter" onPress={toggleModal} />
          </Animated.View>
        </View>
      </Modal>
    </View>
  );
};
