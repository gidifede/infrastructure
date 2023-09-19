DROP TABLE IF EXISTS PRODUCT;
CREATE TABLE `PRODUCT` (
  `id` bigint NOT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS PRODUCT_STATUS;
CREATE TABLE `PRODUCT_STATUS` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint DEFAULT NULL,
  `status` varchar(255) NOT NULL,
  `timestamp` timestamp(6) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `PRODUCT_STATUS_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `PRODUCT` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=61 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS PRODUCT_ACCEPTANCE;
CREATE TABLE `PRODUCT_ACCEPTANCE` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint DEFAULT NULL,
  `timestamp` timestamp(6) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `PRODUCT_ACCEPTANCE_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `PRODUCT` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS PRODUCT_ATTRIBUTE;
CREATE TABLE `PRODUCT_ATTRIBUTE` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `PRODUCT_ATTRIBUTE_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `PRODUCT` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS PRODUCT_DELIVERY;
CREATE TABLE `PRODUCT_DELIVERY` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint DEFAULT NULL,
  `timestamp` timestamp(6) NULL DEFAULT NULL,
  `STATUS` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `PRODUCT_DELIVERY_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `PRODUCT` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS PRODUCT_LIST;
CREATE TABLE `PRODUCT_LIST` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `PRODUCT_LIST_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `PRODUCT` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS PRODUCT_PROCESSING;
CREATE TABLE `PRODUCT_PROCESSING` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint DEFAULT NULL,
  `timestamp` timestamp(6) NULL DEFAULT NULL,
  `STATUS` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `PRODUCT_PROCESSING_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `PRODUCT` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS PRODUCT_TRANSPORT;
CREATE TABLE `PRODUCT_TRANSPORT` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint DEFAULT NULL,
  `timestamp` timestamp(6) NULL DEFAULT NULL,
  `STATUS` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `PRODUCT_TRANSPORT_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `PRODUCT` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS PRODUCT_WITHDRAWAL;
CREATE TABLE `PRODUCT_WITHDRAWAL` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint DEFAULT NULL,
  `timestamp` timestamp(6) NULL DEFAULT NULL,
  `STATUS` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `PRODUCT_WITHDRAWAL_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `PRODUCT` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS ATTRIBUTES_PRODUCT_ACCEPTANCE;
CREATE TABLE `ATTRIBUTES_PRODUCT_ACCEPTANCE` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_acceptance_id` bigint DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_acceptance_id` (`product_acceptance_id`),
  CONSTRAINT `ATTRIBUTES_PRODUCT_ACCEPTANCE_ibfk_1` FOREIGN KEY (`product_acceptance_id`) REFERENCES `PRODUCT_ACCEPTANCE` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS LOCATION;
CREATE TABLE `LOCATION` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint DEFAULT NULL,
  `type` varchar(255) NOT NULL,
  `address` varchar(255) NOT NULL,
  `city` varchar(255) NOT NULL,
  `nation` varchar(255) NOT NULL,
  `zipcode` varchar(255) NOT NULL,
  `location_code` varchar(255) NOT NULL,
  `timestamp` timestamp(6) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `LOCATION_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `PRODUCT` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS LOCATION_ATTRIBUTES;
CREATE TABLE `LOCATION_ATTRIBUTES` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `location_id` bigint DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `location_id` (`location_id`),
  CONSTRAINT `LOCATION_ATTRIBUTES_ibfk_1` FOREIGN KEY (`location_id`) REFERENCES `LOCATION` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS SENDER;
CREATE TABLE `SENDER` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `province` varchar(255) NOT NULL,
  `city` varchar(255) NOT NULL,
  `address` varchar(255) NOT NULL,
  `zipcode` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `SENDER_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `PRODUCT` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS SENDER_ATTRIBUTES;
CREATE TABLE `SENDER_ATTRIBUTES` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `sender_id` bigint DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `sender_id` (`sender_id`),
  CONSTRAINT `SENDER_ATTRIBUTES_ibfk_1` FOREIGN KEY (`sender_id`) REFERENCES `SENDER` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS RECEIVER;
CREATE TABLE `RECEIVER` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `province` varchar(255) NOT NULL,
  `city` varchar(255) NOT NULL,
  `address` varchar(255) NOT NULL,
  `zipcode` varchar(255) NOT NULL,
  `number` int NOT NULL,
  `email` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `RECEIVER_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `PRODUCT` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS RECEIVER_ATTRIBUTES;
CREATE TABLE `RECEIVER_ATTRIBUTES` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `receiver_id` bigint DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `receiver_id` (`receiver_id`),
  CONSTRAINT `RECEIVER_ATTRIBUTES_ibfk_1` FOREIGN KEY (`receiver_id`) REFERENCES `RECEIVER` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;




DROP TABLE IF EXISTS ATTRIBUTES_PRODUCT_DELIVERY;
CREATE TABLE `ATTRIBUTES_PRODUCT_DELIVERY` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_delivery_id` bigint DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_delivery_id` (`product_delivery_id`),
  CONSTRAINT `ATTRIBUTES_PRODUCT_DELIVERY_ibfk_1` FOREIGN KEY (`product_delivery_id`) REFERENCES `PRODUCT_DELIVERY` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS DELIVERY_CARRIER;
CREATE TABLE `DELIVERY_CARRIER` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_delivery_id` bigint DEFAULT NULL,
  `type_id` varchar(255) NOT NULL,
  `vehicle_id` varchar(255) NOT NULL,
  `driver_id` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_delivery_id` (`product_delivery_id`),
  CONSTRAINT `DELIVERY_CARRIER_ibfk_1` FOREIGN KEY (`product_delivery_id`) REFERENCES `PRODUCT_DELIVERY` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



DROP TABLE IF EXISTS ATTRIBUTES_PRODUCT_PROCESSING;
CREATE TABLE `ATTRIBUTES_PRODUCT_PROCESSING` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_processing_id` bigint DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_processing_id` (`product_processing_id`),
  CONSTRAINT `ATTRIBUTES_PRODUCT_PROCESSING_ibfk_1` FOREIGN KEY (`product_processing_id`) REFERENCES `PRODUCT_PROCESSING` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



DROP TABLE IF EXISTS ATTRIBUTES_PRODUCT_WITHDRAWAL;
CREATE TABLE `ATTRIBUTES_PRODUCT_WITHDRAWAL` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_withdrawal_id` bigint DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_withdrawal_id` (`product_withdrawal_id`),
  CONSTRAINT `ATTRIBUTES_PRODUCT_WITHDRAWAL_ibfk_1` FOREIGN KEY (`product_withdrawal_id`) REFERENCES `PRODUCT_WITHDRAWAL` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS WITHDRAWAL_CARRIER;
CREATE TABLE `WITHDRAWAL_CARRIER` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_withdrawal_id` bigint DEFAULT NULL,
  `type_id` int NOT NULL,
  `vehicle_id` int NOT NULL,
  `driver_id` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_withdrawal_id` (`product_withdrawal_id`),
  CONSTRAINT `WITHDRAWAL_CARRIER_ibfk_1` FOREIGN KEY (`product_withdrawal_id`) REFERENCES `PRODUCT_WITHDRAWAL` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;



DROP TABLE IF EXISTS ATTRIBUTES_PRODUCT_TRANSPORT;
CREATE TABLE `ATTRIBUTES_PRODUCT_TRANSPORT` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_transport_id` bigint DEFAULT NULL,
  `from_type` varchar(255) NOT NULL,
  `from_address` varchar(255) NOT NULL,
  `from_zipcode` varchar(255) NOT NULL,
  `from_city` varchar(255) NOT NULL,
  `from_nation` varchar(255) NOT NULL,
  `to_type` varchar(255) NOT NULL,
  `to_address` varchar(255) NOT NULL,
  `to_zipcode` varchar(255) NOT NULL,
  `to_city` varchar(255) NOT NULL,
  `to_nation` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_transport_id` (`product_transport_id`),
  CONSTRAINT `ATTRIBUTES_PRODUCT_TRANSPORT_ibfk_1` FOREIGN KEY (`product_transport_id`) REFERENCES `PRODUCT_TRANSPORT` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS TRANSPORT_CARRIER;
CREATE TABLE `TRANSPORT_CARRIER` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_transport_id` bigint DEFAULT NULL,
  `type_id` int NOT NULL,
  `vehicle_id` int NOT NULL,
  `driver_id` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `product_transport_id` (`product_transport_id`),
  CONSTRAINT `TRANSPORT_CARRIER_ibfk_1` FOREIGN KEY (`product_transport_id`) REFERENCES `PRODUCT_TRANSPORT` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS VEHICLE_STATE;
CREATE TABLE `VEHICLE_STATE` (
   `id` bigint NOT NULL AUTO_INCREMENT,
    `vehicle_id` varchar (255) NOT NULL,
    `route_id` varchar (255) NOT NULL,
    `current_load` int NOT NULL,
    `started` boolean NOT NULL,
    `ended` boolean NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS CENTER;
CREATE TABLE `CENTER` (
  `id` varchar (255) NOT NULL,
  `type` varchar (255) NOT NULL,
  `num_parcels_in_queue_to_be_processed` bigint NOT NULL,
  `num_parcels_ready_to_be_delivered` bigint NOT NULL,
  `num_parcels_undelivered`bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS CENTERS_TRANSPORT;
CREATE TABLE `CENTERS_TRANSPORT` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `source_center_id` varchar (255) NOT NULL,
  `dest_center_id` varchar (255) NOT NULL,
  `num_parcels_in_queue` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS NODE;
CREATE TABLE NODE (
   `id` bigint NOT NULL AUTO_INCREMENT,
   `type` varchar (255) NOT NULL,
   `company` varchar (255) NOT NULL,
   `address` varchar (255) NOT NULL UNIQUE,
    `zipcode` varchar (255) NOT NULL UNIQUE,
    `city` varchar (255) NOT NULL UNIQUE,
    `nation` varchar (255) NOT NULL UNIQUE,
    `latitude` varchar (255) NOT NULL UNIQUE,
    `longitude` varchar (255) NOT NULL UNIQUE,
    PRIMARY KEY (`id`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS SORTING_MACHINE;
CREATE TABLE SORTING_MACHINE (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `node_id` bigint FOREIGN KEY REFERENCES NODE(id),
  `serial` varchar (255) NOT NULL UNIQUE,
  `type` varchar (255) NOT NULL,
  `capacity` varchar (255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS EDGE;
CREATE TABLE EDGE(
   `id` bigint NOT NULL AUTO_INCREMENT,
   `node_source_id` bigint FOREIGN KEY REFERENCES NODE(id),
   `target_source_id` bigint FOREIGN KEY REFERENCES NODE(id),
   `distance` bigint NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS ROUTE;
CREATE TABLE ROUTE(
   `id` bigint NOT NULL AUTO_INCREMENT,
   `edge_id` bigint FOREIGN KEY REFERENCES EDGE(id),
   `cut_off_time` varchar (255) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE USER 'query_handler'@'%' IDENTIFIED BY 'query_handler_password';
GRANT SELECT ON * . * TO 'query_handler'@'%';
CREATE USER 'projection_handler'@'%' IDENTIFIED BY 'projection_handler_password';
GRANT CREATE, SELECT, INSERT, UPDATE ON * . * TO 'projection_handler'@'%';
GRANT DELETE ON CENTER TO query_handler;
GRANT DELETE ON CENTERS_TRANSPORT TO query_handler;
GRANT DELETE ON VEHICLE_STATE TO query_handler;
FLUSH PRIVILEGES;





  

