import subprocess

def kill_processes_on_port(port):
    # Step 1: Find processes listening on the specified port
    try:
        # Run lsof to find the processes listening on the port
        result = subprocess.run(
            ["lsof", "-i", f":{port}"],
            text=True,
            capture_output=True,
            check=True
        )
        
        # Parse the output and collect process IDs (PIDs)
        pids = set()
        for line in result.stdout.splitlines()[1:]:  # Skip the header line
            parts = line.split()
            if len(parts) > 1:
                pids.add(parts[1])

        if not pids:
            print(f"No processes found listening on port {port}.")
            return

        # Step 2: Kill each process by PID
        for pid in pids:
            try:
                subprocess.run(["kill", "-9", pid], check=True)
                print(f"Killed process with PID: {pid}")
            except subprocess.CalledProcessError as e:
                print(f"Failed to kill process with PID {pid}: {e}")

    except subprocess.CalledProcessError as e:
        print(f"Error running lsof: {e}")

# Example usage
kill_processes_on_port(8080)