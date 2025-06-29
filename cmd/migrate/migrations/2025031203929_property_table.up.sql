-- Cities
CREATE TABLE cities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Localities
CREATE TABLE localities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    city_id INT,
    name VARCHAR(100) NOT NULL,
    FOREIGN KEY (city_id) REFERENCES cities(id)
);

-- Property Types (Rent, Resale, Plot, Land)
CREATE TABLE property_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type_name VARCHAR(50) NOT NULL
);

-- Apartment Types
CREATE TABLE apartment_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- BHK Types
CREATE TABLE bhk_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    bhk_label VARCHAR(10) NOT NULL -- e.g., 1BHK, 2BHK
);

-- Facing Directions
CREATE TABLE facing_directions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    direction VARCHAR(20) NOT NULL
);

-- Furnishing Types
CREATE TABLE furnishing_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(50) NOT NULL -- Fully, Semi, Unfurnished
);

-- Preferred Tenants
CREATE TABLE preferred_tenants (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(50) NOT NULL -- Anyone, Family, Bachelor Female, etc.
);

-- Water Supply Types
CREATE TABLE water_supply_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(50) NOT NULL -- Borewell, Corporation
);

-- Property Conditions
CREATE TABLE property_conditions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    condition VARCHAR(100) NOT NULL -- e.g., Excellent, Good, etc.
);

-- Main Properties Table
CREATE TABLE properties (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255),
    description TEXT,
    city_id INT,
    locality_id INT,
    landmark VARCHAR(255),

    property_type_id INT,
    apartment_type_id INT,
    apartment_name VARCHAR(100),
    bhk_type_id INT,
    floor INT,
    total_floors INT,
    property_age VARCHAR(50),
    facing_id INT,
    built_up_area INT,

    bathroom_count INT,
    balcony_count INT,
    parking_available BOOLEAN,

    water_supply_id INT,
    property_condition_id INT,

    available_for ENUM('Rent', 'Lease', 'Sale'),
    expected_rent DECIMAL(10,2),
    expected_deposit DECIMAL(10,2),
    is_rent_negotiable BOOLEAN,
    monthly_maintenance DECIMAL(10,2),
    available_from DATE,

    furnishing_type_id INT,
    contact_number VARCHAR(15),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (city_id) REFERENCES cities(id),
    FOREIGN KEY (locality_id) REFERENCES localities(id),
    FOREIGN KEY (property_type_id) REFERENCES property_types(id),
    FOREIGN KEY (apartment_type_id) REFERENCES apartment_types(id),
    FOREIGN KEY (bhk_type_id) REFERENCES bhk_types(id),
    FOREIGN KEY (facing_id) REFERENCES facing_directions(id),
    FOREIGN KEY (furnishing_type_id) REFERENCES furnishing_types(id),
    FOREIGN KEY (water_supply_id) REFERENCES water_supply_types(id),
    FOREIGN KEY (property_condition_id) REFERENCES property_conditions(id)
);

-- Many-to-Many: Property - Preferred Tenants
CREATE TABLE property_preferred_tenants (
    id INT AUTO_INCREMENT PRIMARY KEY,
    property_id INT,
    tenant_type_id INT,
    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
    FOREIGN KEY (tenant_type_id) REFERENCES preferred_tenants(id)
);
-- Many-to-Many: Property - Images
CREATE TABLE property_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    property_id INT,
    image_url VARCHAR(255) NOT NULL,
    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE
);
-- Many-to-Many: Property - Amenities
CREATE TABLE property_amenities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    property_id INT,
    amenity VARCHAR(100) NOT NULL,
    FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE
);
