package com.example.neer.Sambad.chatRoom;

import java.util.Optional;

import org.springframework.data.mongodb.repository.MongoRepository;

public interface chatRoomRepository extends MongoRepository<chatRoom, String>{
    Optional<chatRoom> findBySenderIdAndRecipientId(String senderId, String recipientId);
}
