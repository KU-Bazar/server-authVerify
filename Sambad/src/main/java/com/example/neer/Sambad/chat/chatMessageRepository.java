package com.example.neer.Sambad.chat;

import org.springframework.data.mongodb.repository.MongoRepository;
import java.util.List;

public interface chatMessageRepository extends MongoRepository<chatMessage, String> {
    List<chatMessage> findByChatId(String chatId);

}
